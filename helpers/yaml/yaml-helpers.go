package yamlHelper

import (
	"fmt"
	"os"
	simplismTypes "simplism/types"

	"gopkg.in/yaml.v3"
)

// ReadYamlFile reads a YAML file and returns a map of server.WasmArguments and an error.
//
// It takes a string parameter called yamlFilePath, which represents the path to the YAML file.
// The function returns a map[string]server.WasmArguments, which is a map of server.WasmArguments objects, and an error.
func ReadYamlFile(yamlFilePath string) (map[string]simplismTypes.WasmArguments, error) {

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

func WriteYamlFile(yamlFilepath string, wasmArgumentsMap map[string]simplismTypes.WasmArguments) {
	yamlFile, err := os.Create(yamlFilepath)
	if err != nil {
		fmt.Println("😡 (writeYamlFile) creating the yaml config file:", err)
		return
	}
	defer yamlFile.Close()

	yamlData, errMarshal := yaml.Marshal(&wasmArgumentsMap)
	if errMarshal != nil {
		fmt.Println("😡 (writeYamlFile) marshalling the yaml config file:", errMarshal)
		return
	}

	_, errWrite := yamlFile.Write(yamlData)
	if errWrite != nil {
		fmt.Println("😡 (writeYamlFile) writing the yaml config file:", errWrite)
		return
	}
}
