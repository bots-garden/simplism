package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"simplism/server"
	"sync"
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

	var serviceDiscoveryWasmArguments server.WasmArguments

	var protection = sync.Mutex{}
	var wasmServices = make(map[string][]string)

	fmt.Println("🐑 flock mode activated")

	ctx := context.Background()

	//ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//defer stop()

	// loop through the map
	for configKey, wasmArguments := range wasmArgumentsMap {
		wasmArguments = applyDefaultValuesIfMissing(wasmArguments)

		if configKey == "service-discovery" {
			serviceDiscoveryWasmArguments = wasmArguments
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
					fmt.Println("😡 Error when starting a new simplism process:", configKey, err)
				} else {
					protection.Lock()
					defer protection.Unlock()
					//TODO: test if certificate to determine if https or not
					wasmServices[configKey] = []string{wasmArguments.HTTPPort, wasmArguments.FunctionName}
				}

			}(configKey, wasmArguments)
		}

	}
	// create a small service-discovery http server
	if serviceDiscoveryWasmArguments.HTTPPort != "" {
		go func() {

			http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

				// return a json list of the wasm plugins
				response.Header().Set("Content-Type", "application/json")
				// transform the map to json
				wasmServicesJSON, err := json.Marshal(wasmServices)
				if err != nil {
					response.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(response, "😡 Error when transforming the wasmServices map to json:", err)
				}
				response.WriteHeader(http.StatusOK)
				fmt.Fprintln(response, string(wasmServicesJSON))
			})

			if serviceDiscoveryWasmArguments.CertFile != "" && serviceDiscoveryWasmArguments.KeyFile != "" {
				// Path to the TLS certificate and key files
				certFile := serviceDiscoveryWasmArguments.CertFile
				keyFile := serviceDiscoveryWasmArguments.KeyFile

				fmt.Println("🔎 http(s) service-discovery flock server is listening on:", serviceDiscoveryWasmArguments.HTTPPort)

				err := http.ListenAndServeTLS(":"+serviceDiscoveryWasmArguments.HTTPPort, certFile, keyFile, nil)
				if err != nil {
					log.Fatal("😡", err)
					os.Exit(1) // 🤔
				} 
			} else {
				fmt.Println("🔎 http service-discovery flock server is listening on:", serviceDiscoveryWasmArguments.HTTPPort)

				err := http.ListenAndServe(":"+serviceDiscoveryWasmArguments.HTTPPort, nil)
				if err != nil {
					log.Fatal("😡", err)
					os.Exit(1) // 🤔
				} 
			}

		}()
	}

	// Listen for the interrupt signal.
	<-ctx.Done()
}
