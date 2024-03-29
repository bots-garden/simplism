package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simplism/server/router"
	simplismTypes "simplism/types"
)

// goRoutineStartServer starts an HTTP server using the provided configuration and arguments.
//
// Parameters:
// - configKey: The key for the configuration. If empty, a default configuration will be used.
// - wasmArgs: The arguments for the WebAssembly module.
//
// Return type: None.
func goRoutineStartServer(configKey string, wasmArgs simplismTypes.WasmArguments) {

	if wasmArgs.CertFile != "" && wasmArgs.KeyFile != "" {
		var message string
		if configKey == "" {
			message = "🌍 http(s) server is listening on: " + wasmArgs.HTTPPort
		} else {
			message = "🌍 [" + configKey + "] http(s) server is listening on: " + wasmArgs.HTTPPort
		}

		// Path to the TLS certificate and key files
		certFile := wasmArgs.CertFile
		keyFile := wasmArgs.KeyFile

		fmt.Println(message)
		err := http.ListenAndServeTLS(":"+wasmArgs.HTTPPort, certFile, keyFile, router.GetRouter())
		//err := http.ListenAndServeTLS(":"+wasmArgs.HTTPPort, certFile, keyFile, nil)

		//currentSimplismProcess.Host =

		if err != nil {
			log.Fatal("😡", err)
			os.Exit(1)
		}
	} else {
		var message string
		if configKey == "" {
			message = "🌍 http server is listening on: " + wasmArgs.HTTPPort
		} else {
			message = "🌍 [" + configKey + "] http(s) server is listening on: " + wasmArgs.HTTPPort
		}
		fmt.Println(message)
		err := http.ListenAndServe(":"+wasmArgs.HTTPPort, router.GetRouter())
		//err := http.ListenAndServe(":"+wasmArgs.HTTPPort, nil)

		if err != nil {
			log.Fatal("😡", err)
			os.Exit(1)
		}
	}
}
