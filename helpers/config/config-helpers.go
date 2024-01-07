package configHelper

import (
	"encoding/json"
	"fmt"
	"os"
	simplismTypes "simplism/types"
)

// GetEnvVarsFromString returns a slice of strings parsed from the given envars string.
//
// Parameters:
// - envars: A string representing the environment variables in JSON format.
//
// Returns:
// - []string: A slice of strings containing the parsed environment variables.
func GetEnvVarsFromString(envars string) []string {
	var vars []string
	unmarshallError := json.Unmarshal([]byte(envars), &vars)
	if unmarshallError != nil {
		fmt.Println("😡 getEnvVarsFromString:", unmarshallError)
		os.Exit(1)
	}
	return vars
}

// GetHostsFromString gets a string representing a JSON array of hosts and returns a slice of strings containing the hosts.
//
// allowHosts: a string representing a JSON array of hosts.
// []string: a slice of strings containing the hosts.
func GetHostsFromString(allowHosts string) []string {
	var hosts []string
	unmarshallError := json.Unmarshal([]byte(allowHosts), &hosts)
	if unmarshallError != nil {
		fmt.Println("😡 getHostsFromString:", unmarshallError)
		os.Exit(1)
	}
	return hosts

}

// GetPathsFromJSONString parses a JSON string and returns a map of paths.
//
// It takes a string parameter `allowPaths` which represents the JSON string to be parsed.
// The function returns a map of type `map[string]string` which contains the parsed paths.
func GetPathsFromJSONString(allowPaths string) map[string]string {
	var paths map[string]string
	unmarshallError := json.Unmarshal([]byte(allowPaths), &paths)
	if unmarshallError != nil {
		fmt.Println("😡 getPathsFromJSONString:", unmarshallError)
		os.Exit(1)
	}
	return paths
}

// GetConfigFromJSONString retrieves a map of configuration properties from a JSON string.
//
// config: a JSON string representing the configuration properties.
// Returns: a map of configuration properties, where the keys are strings and the values are strings.
func GetConfigFromJSONString(config string) map[string]string {
	var manifestConfig map[string]string
	unmarshallError := json.Unmarshal([]byte(config), &manifestConfig)
	if unmarshallError != nil {
		fmt.Println("😡 getConfigFromJSONString:", unmarshallError)
		os.Exit(1)
	}
	return manifestConfig
}

// ApplyDefaultValuesIfMissing applies default values to the given WasmArguments struct if any of its fields are empty.
//
// Parameters:
// - wasmArguments: The WasmArguments struct to apply default values to.
//
// Return type:
// - server.WasmArguments: The WasmArguments struct with default values applied.
func ApplyDefaultValuesIfMissing(wasmArguments simplismTypes.WasmArguments) simplismTypes.WasmArguments {
	// default values:
	if wasmArguments.AllowHosts == "" {
		wasmArguments.AllowHosts = `["*"]`
	}

	if wasmArguments.AllowPaths == "" {
		wasmArguments.AllowPaths = "{}"
	}

	if wasmArguments.Config == "" {
		wasmArguments.Config = "{}"
	}

	if wasmArguments.HTTPPort == "" {
		wasmArguments.HTTPPort = "8080"
	}

	if wasmArguments.EnvVars == "" {
		wasmArguments.EnvVars = "[]"
	}

	if wasmArguments.LogLevel == "" {
		wasmArguments.LogLevel = "info"
	}

	if wasmArguments.RecoveryPath == "" {
		wasmArguments.RecoveryPath = "recovery.yaml"
	}

	return wasmArguments

}
