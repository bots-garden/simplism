package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	wasmHelper "simplism/extism-runtime"
	functionTypes "simplism/function-types"
	httpHelper "simplism/http-helper"

	"github.com/go-resty/resty/v2"
)

// WasmArguments type
type WasmArguments struct {
	FilePath        string
	FunctionName    string
	HTTPPort        string
	Input           string
	LogLevel        string
	AllowHosts      string
	AllowPaths      string
	Config          string
	Wasi            bool
	URL             string
	AuthHeaderName  string
	AuthHeaderValue string
}

func getHostsFromString(allowHosts string) []string {
	var hosts []string
	unmarshallError := json.Unmarshal([]byte(allowHosts), &hosts)
	if unmarshallError != nil {
		fmt.Println(unmarshallError)
		os.Exit(1)
	}
	return hosts

}

func getPathsFromJSONString(allowPaths string) map[string]string {
	var paths map[string]string
	unmarshallError := json.Unmarshal([]byte(allowPaths), &paths)
	if unmarshallError != nil {
		fmt.Println(unmarshallError)
		os.Exit(1)
	}
	return paths
}

func getConfigFromJSONString(config string) map[string]string {
	var manifestConfig map[string]string
	unmarshallError := json.Unmarshal([]byte(config), &manifestConfig)
	if unmarshallError != nil {
		fmt.Println(unmarshallError)
		os.Exit(1)
	}
	return manifestConfig
}

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
		return errors.New("error while downloading the wasm file")
	}

	if err != nil {
		return err
	}
	return nil
}

// Listen wip...
func Listen(wasmArgs WasmArguments) {

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

	level := wasmHelper.GetLevel(wasmArgs.LogLevel)

	ctx := context.Background()

	config, manifest := wasmHelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)

	wasmHelper.GeneratePluginsPool(ctx, config, manifest)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		var (
			result []byte
			err    error
		)

		body := httpHelper.GetBody(request) // is the body the same with fiber ?

		mainFunctionArgument := functionTypes.Argument{
			Header: request.Header,
			Body:   string(body),
			Method: request.Method,
			URI:    request.RequestURI,
		}

		//result, err = wasmHelper.CallWasmFunction(wasmFunctionName, []byte(mainFunctionArgument.ToEncodedJSONString()))
		result, err = wasmHelper.CallWasmFunction(wasmArgs.FunctionName, mainFunctionArgument.ToJSONBuffer())

		/* Expected response
		type ReturnValue struct {
			Body   string              `json:"body"`
			Header map[string][]string `json:"header"`
		}
		*/
		returnValue := functionTypes.ReturnValue{}
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
		// certificate - https ?
		fmt.Println("🌍 http server is listening on:", wasmArgs.HTTPPort)
		log.Fatal(http.ListenAndServe(":"+wasmArgs.HTTPPort, nil))
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

}
