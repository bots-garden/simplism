package cmds

import (
	"fmt"
	"os"
	simplismTypes "simplism/types"

	"gopkg.in/yaml.v3"
)

// readYamlFile reads a YAML file and returns a map of server.WasmArguments and an error.
//
// It takes a string parameter called yamlFilePath, which represents the path to the YAML file.
// The function returns a map[string]server.WasmArguments, which is a map of server.WasmArguments objects, and an error.
func readYamlFile(yamlFilePath string) (map[string]simplismTypes.WasmArguments, error) {

	yamlFile, err := os.ReadFile(yamlFilePath)

	if err != nil {
		return nil, err
	}

	data := make(map[string]simplismTypes.WasmArguments)

	err = yaml.Unmarshal(yamlFile, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}


// getWasmArgumentsMap returns a map of WasmArguments based on the provided configFilepath.
//
// configFilepath: The filepath of the YAML config file.
// map[string]server.WasmArguments: The map of WasmArguments.
func getWasmArgumentsMap(configFilepath string) map[string]simplismTypes.WasmArguments {
	wasmArgumentsMap, err := readYamlFile(configFilepath)
	if err != nil {
		fmt.Println("😡 (getWasmArgumentsMap) reading the yaml config file:", err)
		os.Exit(1)
	}
	//fmt.Println("🤖 +>", wasmArgumentsMap)
	return wasmArgumentsMap
}
