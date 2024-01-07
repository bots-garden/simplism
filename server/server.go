package server

import (
	"context"
	"fmt"

	//"net/http"
	"os"
	"os/signal"
	configHelper "simplism/helpers/config"
	wasmHelper "simplism/helpers/wasm"
	simplismTypes "simplism/types"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"embed"
)

//go:embed embedded
var fs embed.FS

var currentSimplismProcess = simplismTypes.SimplismProcess{}

var router = chi.NewRouter()

// Listen is a function that listens for incoming HTTP requests and processes them using WebAssembly.
//
// It takes a `wasmArgs` parameter of type `WasmArguments` which contains the necessary arguments for configuring the WebAssembly environment.
// The function does not return anything.
func Listen(wasmArgs simplismTypes.WasmArguments, configKey string) {

	// if you don't want to serve a wasm file use "?" instead of the path to the wasm file
	// then at start, Simplism will provide a scratch.wasm file in the root directory
	if wasmArgs.FilePath == "?" && wasmArgs.FunctionName == "?" {
		wasmScratchfile, _ := fs.ReadFile("embedded/scratch.wasm")
		// copy this file to the root directory
		err := os.WriteFile("scratch.wasm", wasmScratchfile, 0644)
		if err != nil {
			fmt.Println("😡 Error copying file to root directory:", err)
			return
		}
		wasmArgs.FilePath = "scratch.wasm"
		wasmArgs.FunctionName = "handle"
	}

	if wasmArgs.FilePath == "?" && wasmArgs.FunctionName != "?" {
		fmt.Println("😡 You have to use ? for the wasm file path and the function name")
		os.Exit(1)
	}


	// Store information about the current simplism process
	currentSimplismProcess.PID = os.Getpid()
	currentSimplismProcess.FilePath = wasmArgs.FilePath
	currentSimplismProcess.FunctionName = wasmArgs.FunctionName
	currentSimplismProcess.HTTPPort = wasmArgs.HTTPPort

	currentSimplismProcess.Information = wasmArgs.Information
	currentSimplismProcess.ServiceName = wasmArgs.ServiceName

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

	//ctx := context.Background()
	// Create context that listens for the interrupt signal from the OS.
	// This context will be used for function calls.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, manifest := wasmHelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)
	wasmHelper.GeneratePluginsPool(ctx, config, manifest)

	/*
		This handler is responsible for:
		- handling HTTP requests and,
		- calling the WebAssembly function.
	*/
	router.HandleFunc("/", mainHandler(wasmArgs))
	//http.HandleFunc("/", mainHandler(wasmArgs))

	// This handler is responsible for reloading the WebAssembly file,
	router.HandleFunc("/reload", reloadHandler(ctx, wasmArgs))
	//http.HandleFunc("/reload", reloadHandler(ctx, wasmArgs))

	// This handler is responsible for listening for the other Simplism processes,
	// The current Simplism process is responsible for handling the list of the other Simplism processes.
	if wasmArgs.ServiceDiscovery == true {
		fmt.Println("🤖 this service is a service discovery")
		router.HandleFunc("/discovery", discoveryHandler(wasmArgs))
		//http.HandleFunc("/discovery", discoveryHandler(wasmArgs))
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
		router.HandleFunc("/spawn", spawnHandler(wasmArgs))

		/* 
			Try to load and start the previous wasm plug-ins
			only if discovery mode is activated 
		*/




	}

	// https://github.com/etcd-io/bbolt
	if wasmArgs.StoreMode == true {
		fmt.Println("📦 this service can store data")
		router.HandleFunc("/store", storeHandler(wasmArgs))
		//http.HandleFunc("/store", storeHandler(wasmArgs))
	}

	// this does not really work
	if wasmArgs.RegistryMode == true {
		fmt.Println("🐳 small wasm registry activated")
		router.HandleFunc("/registry/push", registryHandler(wasmArgs))
		router.HandleFunc("/registry/pull/{wasmfilename}", registryHandler(wasmArgs))

		// TODO: to be implemented in the future 🚧 (soon)
		router.HandleFunc("/registry/remove/{wasmfilename}", registryHandler(wasmArgs))
		router.HandleFunc("/registry/discover", registryHandler(wasmArgs))
		//router.HandleFunc("/registry/discover/{filter}r", registryHandler(wasmArgs))

		//http.HandleFunc("/registry", registryHandler(wasmArgs))
	}

	// Start the Simplism HTTP server
	go func(configKey string) {
		goRoutineStartServer(configKey, wasmArgs)
	}(configKey)

	// Listen for the interrupt signal.
	<-ctx.Done()
	//stop()
	fmt.Println("😢", configKey, "service exited")

}
