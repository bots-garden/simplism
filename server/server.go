package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"simplism/wasmhelper"
	"time"
	simplismTypes "simplism/types"

)

// Listen is a function that listens for incoming HTTP requests and processes them using WebAssembly.
//
// It takes a `wasmArgs` parameter of type `WasmArguments` which contains the necessary arguments for configuring the WebAssembly environment.
// The function does not return anything.
func Listen(wasmArgs simplismTypes.WasmArguments, configKey string) {

	/*
		db, err := bolt.Open("simplism.db", 0600, nil)
		if err != nil {
			log.Fatal("😡🤬🥵", err)
		}
		defer db.Close()
	*/

	// Store information about the current simplism process
	currentSimplismProcess.PID = os.Getpid()
	currentSimplismProcess.FilePath = wasmArgs.FilePath
	currentSimplismProcess.FunctionName = wasmArgs.FunctionName

	currentSimplismProcess.StartTime = time.Now()

	if wasmArgs.URL != "" { // we need to download the wasm file
		fmt.Println("🌍 downloading", wasmArgs.URL, "...")
		err := downloadWasmFile(wasmArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	hosts := getHostsFromString(wasmArgs.AllowHosts)
	paths := getPathsFromJSONString(wasmArgs.AllowPaths)
	manifestConfig := getConfigFromJSONString(wasmArgs.Config)

	// Add environment variable to the manifest config
	envVars := getEnvVarsFromString(wasmArgs.EnvVars)
	// loop throw envVars and add it to the manifest config
	for _, envVar := range envVars {
		manifestConfig[envVar] = os.Getenv(envVar)
	}
	// now we can use `pdk.GetConfig()` to get the value of the environment variables

	level := wasmhelper.GetLevel(wasmArgs.LogLevel)

	ctx := context.Background()

	config, manifest := wasmhelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)

	wasmhelper.GeneratePluginsPool(ctx, config, manifest)

	/*
		This handler is responsible for:
		- handling HTTP requests and,
		- calling the WebAssembly function.
	*/
	http.HandleFunc("/", mainHandler(wasmArgs))

	/*
		This handler is responsible for:
		- reloading the WebAssembly file,
	*/
	http.HandleFunc("/reload", reloadHandler(ctx, wasmArgs))

	/*
		The current Simplism process is responsible for handling the list of the other Simplism processes.
	*/
	if wasmArgs.Discovery == true {
		/*
			This handler is responsible for:
			- listening for the other Simplism processes,
		*/
		http.HandleFunc("/discovery", discoveryHandler(wasmArgs))

		// TODO: add a go routine to clean the process db: remove all old processes
		// Or do it into the discovery handler ?
	}

	/*
		Every 20 seconds, send information about the current simplism process to the discovery simplism process.
	*/
	if wasmArgs.DiscoveryEndpoint != "" {
		fmt.Println("👋 this service is discoverable")
		
		go func() {
			goRoutineSimplismProcess(currentSimplismProcess, wasmArgs)
		}()
	}

	// Start the Simplism HTTP server
	go func(configKey string) {
		goRoutineStartServer(configKey, wasmArgs)
	}(configKey)

	// Listen for the interrupt signal.
	<-ctx.Done()

}
