package cmds

import (
	"fmt"
	"os"
	"simplism/server"

	"gopkg.in/yaml.v3"
)

// readYamlFile reads a YAML file and returns a map of server.WasmArguments and an error.
//
// It takes a string parameter called yamlFilePath, which represents the path to the YAML file.
// The function returns a map[string]server.WasmArguments, which is a map of server.WasmArguments objects, and an error.
func readYamlFile(yamlFilePath string) (map[string]server.WasmArguments, error) {

	yamlFile, err := os.ReadFile(yamlFilePath)

	if err != nil {
		return nil, err
	}

	data := make(map[string]server.WasmArguments)

	err = yaml.Unmarshal(yamlFile, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

// applyDefaultValuesIfMissing applies default values to the given WasmArguments struct if any of its fields are empty.
//
// Parameters:
// - wasmArguments: The WasmArguments struct to apply default values to.
//
// Return type:
// - server.WasmArguments: The WasmArguments struct with default values applied.
func applyDefaultValuesIfMissing(wasmArguments server.WasmArguments) server.WasmArguments {
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

    //TODO: add log level
	return wasmArguments

}

// getWasmArgumentsMap returns a map of WasmArguments based on the provided configFilepath.
//
// configFilepath: The filepath of the YAML config file.
// map[string]server.WasmArguments: The map of WasmArguments.
func getWasmArgumentsMap(configFilepath string) map[string]server.WasmArguments {
	wasmArgumentsMap, err := readYamlFile(configFilepath)
	if err != nil {
		fmt.Println("😡 (getWasmArgumentsMap) reading the yaml config file:", err)
		os.Exit(1)
	}
    return wasmArgumentsMap
}
