package cmds

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"simplism/server"
)

// getExecutablePath returns the path of the executable file for the given program name.
//
// It takes a string parameter `progName` which represents the name of the program.
// It returns a string which represents the path of the executable file.
func getExecutablePath(progName string) string {
	executablePath, err := exec.LookPath("simplism")
	if err != nil {
		fmt.Println("😡 Error finding executable:", err)
		os.Exit(1)
	}
	return executablePath
}

// startFlockMode activates flock mode.
//
// It reads the yaml file to get the wasm arguments of each wasm service. It gets the executable path of simplism.
// It creates a small admin http server. It listens for the interrupt signal.
func startFlockMode(configFilepath string) {

	// read the yaml file to get the wasm arguments of each wasm service
	wasmArgumentsMap := getWasmArgumentsMap(configFilepath)
	// get the executable path of simplism
	simplismExecutablePath := getExecutablePath("simplism")

	//var serviceDiscoveryWasmArguments server.WasmArguments

	fmt.Println("🐑 flock mode activated")

	ctx := context.Background()

	//ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//defer stop()

	// loop through the map
	for configKey, wasmArguments := range wasmArgumentsMap {
		wasmArguments = applyDefaultValuesIfMissing(wasmArguments)

		// Start a new server process with the specified wasm plugin in the config
		go func(configKey string, wasmArguments server.WasmArguments) {
			cmd := &exec.Cmd{
				Path:   simplismExecutablePath,
				Args:   []string{"", "config", configFilepath, configKey},
				Stdout: os.Stdout,
				Stderr: os.Stdout,
			}
			err := cmd.Start()
			if err != nil {
				fmt.Println("😡 Error when starting a new simplism process:", configKey, err)
				os.Exit(1) // exit with an error
			}

		}(configKey, wasmArguments)

	}

	// Listen for the interrupt signal.
	<-ctx.Done()
}
