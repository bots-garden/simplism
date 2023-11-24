package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	functiontypes "simplism/functiontypes"
	httphelper "simplism/httphelper"
	"simplism/wasmhelper"

	"github.com/go-resty/resty/v2"
)

// WasmArguments type
type WasmArguments struct {
	FilePath        string `yaml:"wasm-file,omitempty"`
	FunctionName    string `yaml:"wasm-function,omitempty"`
	HTTPPort        string `yaml:"http-port,omitempty"`
	Input           string `yaml:"input,omitempty"`
	LogLevel        string `yaml:"log-level,omitempty"`
	AllowHosts      string `yaml:"allow-hosts,omitempty"`
	AllowPaths      string `yaml:"allow-paths,omitempty"`
	Config          string `yaml:"config,omitempty"`
	Wasi            bool   `yaml:"wasi,omitempty"`
	URL             string `yaml:"wasm-url,omitempty"`
	AuthHeaderName  string `yaml:"auth-header-name,omitempty"`
	AuthHeaderValue string `yaml:"auth-header-value,omitempty"`
	CertFile        string `yaml:"cert-file,omitempty"`
	KeyFile         string `yaml:"key-file,omitempty"`
}

// getHostsFromString gets a string representing a JSON array of hosts and returns a slice of strings containing the hosts.
//
// allowHosts: a string representing a JSON array of hosts.
// []string: a slice of strings containing the hosts.
func getHostsFromString(allowHosts string) []string {
	var hosts []string
	unmarshallError := json.Unmarshal([]byte(allowHosts), &hosts)
	if unmarshallError != nil {
		fmt.Println("🔴 getHostsFromString:", unmarshallError)
		os.Exit(1)
	}
	return hosts

}

// getPathsFromJSONString parses a JSON string and returns a map of paths.
//
// It takes a string parameter `allowPaths` which represents the JSON string to be parsed.
// The function returns a map of type `map[string]string` which contains the parsed paths.
func getPathsFromJSONString(allowPaths string) map[string]string {
	var paths map[string]string
	unmarshallError := json.Unmarshal([]byte(allowPaths), &paths)
	if unmarshallError != nil {
		fmt.Println("🔴 getPathsFromJSONString:", unmarshallError)
		os.Exit(1)
	}
	return paths
}

// getConfigFromJSONString retrieves a map of configuration properties from a JSON string.
//
// config: a JSON string representing the configuration properties.
// Returns: a map of configuration properties, where the keys are strings and the values are strings.
func getConfigFromJSONString(config string) map[string]string {
	var manifestConfig map[string]string
	unmarshallError := json.Unmarshal([]byte(config), &manifestConfig)
	if unmarshallError != nil {
		fmt.Println("🔴 getConfigFromJSONString:", unmarshallError)
		os.Exit(1)
	}
	return manifestConfig
}

// downloadWasmFile downloads a WebAssembly (Wasm) file from a given URL and saves it to the specified file path.
//
// It takes a WasmArguments struct as a parameter, which contains the necessary information for the download, such as the URL, authentication header, and file path.
// The WasmArguments struct has the following fields:
// - AuthHeaderName (string): the name of the authentication header (e.g., "PRIVATE-TOKEN")
// - AuthHeaderValue (string): the value of the authentication header (e.g., "${GITLAB_WASM_TOKEN}")
// - FilePath (string): the file path where the downloaded Wasm file will be saved
// - URL (string): the URL from which the Wasm file will be downloaded
//
// This function returns an error if there is any issue during the download process, such as a network error or an error response from the server.
// If the download is successful, it returns nil.
func downloadWasmFile(wasmArgs WasmArguments) error {
	// authenticationHeader:
	// Example: "PRIVATE-TOKEN: ${GITLAB_WASM_TOKEN}"
	client := resty.New()

	if wasmArgs.AuthHeaderName != "" {
		client.SetHeader(wasmArgs.AuthHeaderName, wasmArgs.AuthHeaderValue)
	}

	resp, err := client.R().
		SetOutput(wasmArgs.FilePath).
		Get(wasmArgs.URL)

	if resp.IsError() {
		return errors.New("🔴 error while downloading the wasm file")
	}

	if err != nil {
		return err
	}
	return nil
}

// Listen is a function that listens for incoming HTTP requests and processes them using WebAssembly.
//
// It takes a `wasmArgs` parameter of type `WasmArguments` which contains the necessary arguments for configuring the WebAssembly environment.
// The function does not return anything.
func Listen(wasmArgs WasmArguments) {

	// fmt.Println("🤖", wasmArgs)

	if wasmArgs.URL != "" { // we need to download the wasm file
		fmt.Println("🌍 downloading...", wasmArgs.URL)
		err := downloadWasmFile(wasmArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	hosts := getHostsFromString(wasmArgs.AllowHosts)
	paths := getPathsFromJSONString(wasmArgs.AllowPaths)
	manifestConfig := getConfigFromJSONString(wasmArgs.Config)

	level := wasmhelper.GetLevel(wasmArgs.LogLevel)

	ctx := context.Background()

	config, manifest := wasmhelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)

	wasmhelper.GeneratePluginsPool(ctx, config, manifest)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		var (
			result []byte
			err    error
		)

		body := httphelper.GetBody(request) // is the body the same with fiber ?

		mainFunctionArgument := functiontypes.Argument{
			Header: request.Header,
			Body:   string(body),
			Method: request.Method,
			URI:    request.RequestURI,
		}

		//result, err = wasmHelper.CallWasmFunction(wasmFunctionName, []byte(mainFunctionArgument.ToEncodedJSONString()))
		result, err = wasmhelper.CallWasmFunction(wasmArgs.FunctionName, mainFunctionArgument.ToJSONBuffer())

		/* Expected response
		type ReturnValue struct {
			Body   string              `json:"body"`
			Header map[string][]string `json:"header"`
		}
		*/
		returnValue := functiontypes.ReturnValue{}
		errJSONUnmarshal := json.Unmarshal(result, &returnValue)

		if errJSONUnmarshal != nil {
			// send response http code error
			response.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(response, errJSONUnmarshal.Error())
		} else {
			for key, value := range returnValue.Header {
				response.Header().Set(key, value[0])
			}

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(response, err.Error())
			} else {
				// TODO: add default response code if empty?
				response.WriteHeader(returnValue.Code)
				fmt.Fprintln(response, string(returnValue.Body))
			}
		}

	})

	go func() {

		if wasmArgs.CertFile != "" && wasmArgs.KeyFile != "" {
			fmt.Println("🌍 http(s) server is listening on:", wasmArgs.HTTPPort)
			// Path to the TLS certificate and key files
			certFile := wasmArgs.CertFile
			keyFile := wasmArgs.KeyFile

			err := http.ListenAndServeTLS(":"+wasmArgs.HTTPPort, certFile, keyFile, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("🌍 http server is listening on:", wasmArgs.HTTPPort)
			err := http.ListenAndServe(":"+wasmArgs.HTTPPort, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

}
