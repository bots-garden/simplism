package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	configHelper "simplism/helpers/config"
	wasmHelper "simplism/helpers/wasm"
	simplismTypes "simplism/types"
	"time"
)

var currentSimplismProcess = simplismTypes.SimplismProcess{}

// Listen is a function that listens for incoming HTTP requests and processes them using WebAssembly.
//
// It takes a `wasmArgs` parameter of type `WasmArguments` which contains the necessary arguments for configuring the WebAssembly environment.
// The function does not return anything.
func Listen(wasmArgs simplismTypes.WasmArguments, configKey string) {


	// Store information about the current simplism process
	currentSimplismProcess.PID = os.Getpid()
	currentSimplismProcess.FilePath = wasmArgs.FilePath
	currentSimplismProcess.FunctionName = wasmArgs.FunctionName
	currentSimplismProcess.HTTPPort = wasmArgs.HTTPPort

	currentSimplismProcess.Information = wasmArgs.Information
	currentSimplismProcess.ServiceName = wasmArgs.ServiceName

	// for debugging
	//fmt.Println("🤔", currentSimplismProcess.Information,":", currentSimplismProcess.ServiceName)

	currentSimplismProcess.StartTime = time.Now()

	if wasmArgs.URL != "" { // we need to download the wasm file
		fmt.Println("🌍 downloading", wasmArgs.URL, "...")
		err := wasmHelper.DownloadWasmFile(wasmArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	hosts := configHelper.GetHostsFromString(wasmArgs.AllowHosts)
	paths := configHelper.GetPathsFromJSONString(wasmArgs.AllowPaths)
	manifestConfig := configHelper.GetConfigFromJSONString(wasmArgs.Config)

	// Add environment variable to the manifest config
	envVars := configHelper.GetEnvVarsFromString(wasmArgs.EnvVars)
	// loop throw envVars and add it to the manifest config
	for _, envVar := range envVars {
		manifestConfig[envVar] = os.Getenv(envVar)
	}
	// now we can use `pdk.GetConfig()` to get the value of the environment variables

	level := wasmHelper.GetLevel(wasmArgs.LogLevel)

	ctx := context.Background()

	config, manifest := wasmHelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)

	wasmHelper.GeneratePluginsPool(ctx, config, manifest)

	/*
	This handler is responsible for:
	- handling HTTP requests and,
	- calling the WebAssembly function.
	*/
	http.HandleFunc("/", mainHandler(wasmArgs))

	// This handler is responsible for reloading the WebAssembly file,
	http.HandleFunc("/reload", reloadHandler(ctx, wasmArgs))


	// This handler is responsible for listening for the other Simplism processes,
	// The current Simplism process is responsible for handling the list of the other Simplism processes.
	if wasmArgs.ServiceDiscovery == true {
		fmt.Println("🤖 this service is a service discovery")
		http.HandleFunc("/discovery", discoveryHandler(wasmArgs))
	}

	//Every N seconds, send information about the current simplism process to the discovery simplism process.
	// TODO: add a paramater for this
	if wasmArgs.DiscoveryEndpoint != "" {
		fmt.Println("👋 this service is discoverable")
		go func() {
			goRoutineSimplismProcess(currentSimplismProcess, wasmArgs)
		}()
	}

	// This handler is responsible for spawning other services
	// That means that the current simplism process can spawn other simplism processes
	if wasmArgs.SpawnMode == true {
		fmt.Println("🚀 this service can spawn other services")
		http.HandleFunc("/spawn", spawnHandler(wasmArgs))
	}

	// Start the Simplism HTTP server
	go func(configKey string) {
		goRoutineStartServer(configKey, wasmArgs)
	}(configKey)

	// Listen for the interrupt signal.
	<-ctx.Done()

}
