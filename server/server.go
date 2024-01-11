package server

import (
	"context"
	"fmt"

	"os"
	"os/signal"

	configHelper "simplism/helpers/config"
	wasmHelper "simplism/helpers/wasm"
	yamlHelper "simplism/helpers/yaml"
	simplismTypes "simplism/types"

	"syscall"
	"time"

	"embed"

	"simplism/server/discovery"
	"simplism/server/processes"
	"simplism/server/reload"
	wasmfunction "simplism/server/wasm-function"

	"simplism/server/registry"
	"simplism/server/router"
	"simplism/server/spawn"
	"simplism/server/store"
)

//go:embed embedded
var fs embed.FS

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

	/*
	processes.SetCurrentProcessPID(os.Getpid())
	processes.SetCurrentProcessFilePath(wasmArgs.FilePath)
	processes.SetCurrentProcessFunctionName(wasmArgs.FunctionName)
	processes.SetCurrentProcessHTTPPort(wasmArgs.HTTPPort)
	processes.SetCurrentProcessInformation(wasmArgs.Information)
	processes.SetCurrentProcessServiceName(wasmArgs.ServiceName)
	processes.SetCurrentProcessStartTime(time.Now())
	*/

	processes.SetCurrentProcessData(simplismTypes.SimplismProcess{
		PID: os.Getpid(),
		FilePath: wasmArgs.FilePath,
		FunctionName: wasmArgs.FunctionName,
		HTTPPort: wasmArgs.HTTPPort,
		Information: wasmArgs.Information,
		ServiceName: wasmArgs.ServiceName,
		StartTime: time.Now(),
	})


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
	router.GetRouter().HandleFunc("/", wasmfunction.Handler(wasmArgs))

	// This handler is responsible for reloading the WebAssembly file,
	router.GetRouter().HandleFunc("/reload", reload.Handler(ctx, wasmArgs))

	// This handler is responsible for listening for the other Simplism processes,
	// The current Simplism process is responsible for handling the list of the other Simplism processes.
	if wasmArgs.ServiceDiscovery == true {
		fmt.Println("🤖 this service is a service discovery")

		// Delete the db file
		err := os.Remove(wasmArgs.FilePath + ".processes.db")
		if err != nil {
			fmt.Println("😡 Error deleting the db file:", err)
		}

		router.GetRouter().HandleFunc("/discovery", discovery.Handler(wasmArgs))
		//router.GetRouter().HandleFunc("/discovery/{option}", discovery.Handler(wasmArgs))

	}

	//Every N seconds, send information about the current simplism process to the discovery simplism process.
	// TODO: add a paramater for this
	if wasmArgs.DiscoveryEndpoint != "" {
		fmt.Println("👋 this service is discoverable")
		go func() {
			goRoutineSimplismProcess(processes.GetCurrentSimplismProcess(), wasmArgs)
		}()
	}

	// This handler is responsible for spawning other services
	// That means that the current simplism process can spawn other simplism processes
	if wasmArgs.SpawnMode == true {

		fmt.Println("🚀 this service can spawn other services")
		router.GetRouter().HandleFunc("/spawn", spawn.Handler(wasmArgs))

		// TODO: check if a recovery file is existing
		// Read the recovery file and rename it
		if wasmArgs.RecoveryMode == true {

			fmt.Println("🛟 recovery mode activated", wasmArgs.RecoveryPath)

			formerProcessesArguments, err := yamlHelper.ReadYamlFile(wasmArgs.RecoveryPath)
			if err == nil {
				spawn.NotifySpawnServiceForRecovery(formerProcessesArguments)
				// then delete the recovery file ?
				// no because the map of the current running processes is empty at start
				// so the content of the recovery file will be erased anyway
			} else {
				fmt.Println("😡 reading the recovery file:", err)
			}
		}
	}

	// https://github.com/etcd-io/bbolt
	if wasmArgs.StoreMode == true {
		fmt.Println("📦 this service can store data")
		router.GetRouter().HandleFunc("/store", store.Handler(wasmArgs))
		//http.HandleFunc("/store", storeHandler(wasmArgs))
	}

	// this is a 🚧 WIP
	if wasmArgs.RegistryMode == true {
		fmt.Println("🐳 small wasm registry activated")
		router.GetRouter().HandleFunc("/registry/push", registry.Handler(wasmArgs))
		router.GetRouter().HandleFunc("/registry/pull/{wasmfilename}", registry.Handler(wasmArgs))
		router.GetRouter().HandleFunc("/registry/discover", registry.Handler(wasmArgs))

		// TODO: to be implemented in the future 🚧 (soon)
		router.GetRouter().HandleFunc("/registry/remove/{wasmfilename}", registry.Handler(wasmArgs))
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
