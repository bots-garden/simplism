package cmds

import (
	"flag"
	"fmt"
	"os"
	"simplism/server"
	configHelper "simplism/helpers/config"
)

// startConfigMode is a function that starts the configuration mode.
//
// It takes a configFilepath parameter of type string.
// It does not return any value.
func startConfigMode(configFilepath string) {

	wasmArgumentsMap := getWasmArgumentsMap(configFilepath)

	if len(flag.Args()) <= 2 {
		fmt.Println("😡 you must provide a configuration key")
		os.Exit(1)

	} else {
		configKey := flag.Args()[2]

		// Start the server with the specified wasm plugin in the config
		wasmArguments := wasmArgumentsMap[configKey]
		wasmArguments = configHelper.ApplyDefaultValuesIfMissing(wasmArguments)

		//fmt.Println("🤖🤫 >>", wasmArguments.AdminDiscoveryToken)

		server.Listen(wasmArguments, configKey)
	}
}
