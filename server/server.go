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
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// WasmArguments type
type WasmArguments struct {
	FilePath          string `yaml:"wasm-file,omitempty"`
	FunctionName      string `yaml:"wasm-function,omitempty"`
	HTTPPort          string `yaml:"http-port,omitempty"`
	Input             string `yaml:"input,omitempty"`
	LogLevel          string `yaml:"log-level,omitempty"`
	AllowHosts        string `yaml:"allow-hosts,omitempty"`
	AllowPaths        string `yaml:"allow-paths,omitempty"`
	EnvVars           string `yaml:"env,omitempty"`
	Config            string `yaml:"config,omitempty"`
	Wasi              bool   `yaml:"wasi,omitempty"`
	URL               string `yaml:"wasm-url,omitempty"`
	WasmURLAuthHeader string `yaml:"wasm-url-auth-header,omitempty"`
	//AuthHeaderName  string `yaml:"auth-header-name,omitempty"`
	//AuthHeaderValue string `yaml:"auth-header-value,omitempty"`
	CertFile         string `yaml:"cert-file,omitempty"`
	KeyFile          string `yaml:"key-file,omitempty"`
	AdminReloadToken string `yaml:"admin-reload-token,omitempty"`
}

type SimplismProcess struct {
	PID          int
	FunctionName string
	FilePath     string
}

var currentSimplismProcess = SimplismProcess{}

func getEnvVarsFromString(envars string) []string {
	var vars []string
	unmarshallError := json.Unmarshal([]byte(envars), &vars)
	if unmarshallError != nil {
		fmt.Println("😡 getEnvVarsFromString:", unmarshallError)
		os.Exit(1)
	}
	return vars
}

// getHostsFromString gets a string representing a JSON array of hosts and returns a slice of strings containing the hosts.
//
// allowHosts: a string representing a JSON array of hosts.
// []string: a slice of strings containing the hosts.
func getHostsFromString(allowHosts string) []string {
	var hosts []string
	unmarshallError := json.Unmarshal([]byte(allowHosts), &hosts)
	if unmarshallError != nil {
		fmt.Println("😡 getHostsFromString:", unmarshallError)
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
		fmt.Println("😡 getPathsFromJSONString:", unmarshallError)
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
		fmt.Println("😡 getConfigFromJSONString:", unmarshallError)
		os.Exit(1)
	}
	return manifestConfig
}

func getHeaderFromString(headerNameAndValue string) (string, string) {
	splitHeader := strings.Split(headerNameAndValue, "=")
	headerName := splitHeader[0]
	// join all item of splitAuthHeader with "" except the first one
	headerValue := strings.Join(splitHeader[1:], "")
	return headerName, headerValue
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

	//fmt.Println("🚧 downloading", wasmArgs.FilePath, "...")

	if wasmArgs.WasmURLAuthHeader != "" {
		authHeaderName, authHeaderValue := getHeaderFromString(wasmArgs.WasmURLAuthHeader)
		client.SetHeader(authHeaderName, authHeaderValue)

	} else {
		// check if the environment variable WASM_URL_AUTH_HEADER is set
		wasmURLAuthHeader := os.Getenv("WASM_URL_AUTH_HEADER")
		if wasmURLAuthHeader != "" {
			authHeaderName, authHeaderValue := getHeaderFromString(wasmURLAuthHeader)
			client.SetHeader(authHeaderName, authHeaderValue)

		}
	}

	resp, err := client.R().
		SetOutput(wasmArgs.FilePath).
		Get(wasmArgs.URL)

	if resp.IsError() {
		return errors.New("😡 error while downloading the wasm file")
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
func Listen(wasmArgs WasmArguments, configKey string) {

	// fmt.Println("🤖", wasmArgs)
	//fmt.Println("🤖", os.Getpid())

	// Store information about the current simplism process
	currentSimplismProcess.PID = os.Getpid()
	currentSimplismProcess.FilePath = wasmArgs.FilePath
	currentSimplismProcess.FunctionName = wasmArgs.FunctionName

	if wasmArgs.URL != "" { // we need to download the wasm file
		fmt.Println("🌍 downloading", wasmArgs.URL, "...")
		err := downloadWasmFile(wasmArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	hosts := getHostsFromString(wasmArgs.AllowHosts)
	paths := getPathsFromJSONString(wasmArgs.AllowPaths)
	manifestConfig := getConfigFromJSONString(wasmArgs.Config)

	// Add environment variable to the manifest config
	envVars := getEnvVarsFromString(wasmArgs.EnvVars)
	// loop throw envVars and add it to the manifest config
	for _, envVar := range envVars {
		manifestConfig[envVar] = os.Getenv(envVar)
	}
	// now we can use `pdk.GetConfig()` to get the value of the environment variables

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

	http.HandleFunc("/reload", func(response http.ResponseWriter, request *http.Request) {
		// wait for:
		// - POST request
		// - url for the wasm file
		// - admin-reload-token and ADMIN_RELOAD_TOKEN env variable

		// read the header admin-reload-token
		adminReloadToken := request.Header.Get("admin-reload-token")
		var authorised bool = false

		if wasmArgs.AdminReloadToken != "" {
			// token is awaited
			if wasmArgs.AdminReloadToken == adminReloadToken {
				authorised = true
			} else {
				// send response http code error
				response.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(response, "😡 wrong token")
			}

		} else {
			// check if the environment variable WASM_URL_AUTH_HEADER is set
			wasmURLAuthHeader := os.Getenv("ADMIN_RELOAD_TOKEN")
			if wasmURLAuthHeader != "" {
				// token is awaited
				if wasmURLAuthHeader == adminReloadToken {
					authorised = true
				} else {
					// send response http code error
					response.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(response, "😡 wrong token")
				}
			} else {
				authorised = true
			}
		}

		// Test if it's a POST request
		if request.Method == "POST" && authorised == true {
			body := httphelper.GetBody(request)
			// body is a JSON string, extract the url field value of the JSON string
			bodyMap := map[string]string{}
			err := json.Unmarshal([]byte(body), &bodyMap)

			if err != nil {
				// send response http code error
				response.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(response, "😡 "+err.Error())
			} else {

				// TODO test token

				wasmArgs.URL = bodyMap["wasm-url"]
				wasmArgs.FilePath = bodyMap["wasm-file"]
				wasmArgs.FunctionName = bodyMap["wasm-function"]

				fmt.Println("🚀 downloading", wasmArgs.URL, "...")
				err := downloadWasmFile(wasmArgs)
				if err != nil {
					fmt.Println(err)
					os.Exit(1) // TODO: do something else
				}

				hosts := getHostsFromString(wasmArgs.AllowHosts)
				paths := getPathsFromJSONString(wasmArgs.AllowPaths)
				manifestConfig := getConfigFromJSONString(wasmArgs.Config)

				// Add environment variable to the manifest config
				envVars := getEnvVarsFromString(wasmArgs.EnvVars)
				// loop throw envVars and add it to the manifest config
				for _, envVar := range envVars {
					manifestConfig[envVar] = os.Getenv(envVar)
				}
				// now we can use `pdk.GetConfig()` to get the value of the environment variables

				level := wasmhelper.GetLevel(wasmArgs.LogLevel)

				//ctx := context.Background()

				config, manifest := wasmhelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)

				wasmhelper.ReplacePluginInPool(0, ctx, config, manifest)
				fmt.Println("🙂 new wasm plug-in reloaded")

				response.WriteHeader(http.StatusOK)
				fmt.Fprintln(response, string("🙂 new wasm plug-in reloaded"))

				// Update information about the current simplism process
				currentSimplismProcess.FilePath = wasmArgs.FilePath
				currentSimplismProcess.FunctionName = wasmArgs.FunctionName

			}

		} else {
			// response that it's not allowed
			response.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(response, "😡 Method not allowed or you're not authorized")
		}

	})

	/*
		Every 20 seconds, store information about the current simplism process
	*/
	go func() {
		
		// TODO: store the process a first time

		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// TODO: this is a work in progress 🚧
				fmt.Println("👋", currentSimplismProcess)
				// store data somewhere
				// how to garden the data? (with which condition)
				// create an API endpoint to query the data
			}
		}

	}()

	go func(configKey string) {

		if wasmArgs.CertFile != "" && wasmArgs.KeyFile != "" {
			var message string
			if configKey == "" {
				message = "🌍 http(s) server is listening on: " + wasmArgs.HTTPPort
			} else {
				message = "🌍 [" + configKey + "] http(s) server is listening on: " + wasmArgs.HTTPPort
			}

			// Path to the TLS certificate and key files
			certFile := wasmArgs.CertFile
			keyFile := wasmArgs.KeyFile

			fmt.Println(message)
			err := http.ListenAndServeTLS(":"+wasmArgs.HTTPPort, certFile, keyFile, nil)
			if err != nil {
				log.Fatal("😡", err)
				os.Exit(1)
			}
		} else {
			var message string
			if configKey == "" {
				message = "🌍 http server is listening on: " + wasmArgs.HTTPPort
			} else {
				message = "🌍 [" + configKey + "] http(s) server is listening on: " + wasmArgs.HTTPPort
			}
			fmt.Println(message)
			err := http.ListenAndServe(":"+wasmArgs.HTTPPort, nil)
			if err != nil {
				log.Fatal("😡", err)
				os.Exit(1)
			}
		}
	}(configKey)

	// Listen for the interrupt signal.
	<-ctx.Done()

}
