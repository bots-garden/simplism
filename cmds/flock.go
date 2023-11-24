package cmds

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		fmt.Println("🔴 Error finding executable:", err)
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

	var adminWasmArguments server.WasmArguments

	var wasmServices = map[string]string{}

	fmt.Println("🐑 flock mode activated")

	ctx := context.Background()
	// loop through the map
	for configKey, wasmArguments := range wasmArgumentsMap {
		wasmArguments = applyDefaultValuesIfMissing(wasmArguments)

		if configKey == "admin" {
			adminWasmArguments = wasmArguments
		} else {
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
					fmt.Println("🔴 Error when starting a new simplism process:", configKey, err)
				} else {
					wasmServices[configKey] = wasmArguments.HTTPPort
				}

			}(configKey, wasmArguments)
		}

	}
	// create a small admin http server
	if adminWasmArguments.HTTPPort != "" {
		go func() {

			http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
				// TODO: add the list of the wasm plugins
                // TODO: return json list of the wasm plugins + headers
				fmt.Fprintln(response, "🐏 flock mode activated")
				fmt.Fprintln(response, wasmServices)
			})

			if adminWasmArguments.CertFile != "" && adminWasmArguments.KeyFile != "" {
				fmt.Println("🔐 http(s) admin flock server is listening on:", adminWasmArguments.HTTPPort)
				// Path to the TLS certificate and key files
				certFile := adminWasmArguments.CertFile
				keyFile := adminWasmArguments.KeyFile

				err := http.ListenAndServeTLS(":"+adminWasmArguments.HTTPPort, certFile, keyFile, nil)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Println("🔐 http admin flock server is listening on:", adminWasmArguments.HTTPPort)
				err := http.ListenAndServe(":"+adminWasmArguments.HTTPPort, nil)
				if err != nil {
					log.Fatal(err)
				}
			}

		}()
	}

	// Listen for the interrupt signal.
	<-ctx.Done()
}
