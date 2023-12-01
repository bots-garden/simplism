package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

// getEnvVarsFromString returns a slice of strings parsed from the given envars string.
//
// Parameters:
// - envars: A string representing the environment variables in JSON format.
//
// Returns:
// - []string: A slice of strings containing the parsed environment variables.
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

// getHeaderFromString returns the header name and value from a given header string.
//
// The parameter headerNameAndValue is a string that contains the header name and value separated by "=".
// The function splits the headerNameAndValue string using "=" as the separator and stores the result in the splitHeader variable.
// The header name is extracted from the first element of the splitHeader array and stored in the headerName variable.
// The header value is obtained by joining all the elements of the splitHeader array, except the first one, using an empty string as the separator.
// The function then returns the headerName and headerValue as a tuple.
//
// The function returns two string values: headerName and headerValue.
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
