package cmds

import (
	"flag"
	"simplism/server"
)

// startListening initializes a server to listen for requests and executes a WebAssembly function.
//
// Parameters:
// - wasmFilePath: The path to the WebAssembly file.
// - wasmFunctionName: The name of the WebAssembly function to execute.
// - flagSet: The flag set containing the command-line arguments.
// - args: The command-line arguments.
func startListening(wasmFilePath, wasmFunctionName string, flagSet *flag.FlagSet, args []string) {

	httpPort := flagSet.String("http-port", "8080", "http port")

	input := flagSet.String("input", "", "Argument of the function")
	logLevel := flagSet.String("log-level", "", "Log level to print message")
	allowHosts := flagSet.String("allow-hosts", `["*"]`, "Hosts for HTTP request (json array)")
	allowPaths := flagSet.String("allow-paths", "{}", "Allowed paths to write and read files (json string)")
	config := flagSet.String("config", "{}", "Configuration data (json string)")
	wasi := flagSet.Bool("wasi", true, "")

	wasmURL := flagSet.String("wasm-url", "", "Url to download the wasm file")
	authHeaderName := flagSet.String("auth-header-name", "", "Authentication header name, ex: PRIVATE-TOKEN")
	authHeaderValue := flagSet.String("auth-header-value", "", "Value of the authentication header, ex: IlovePandas")

	certFile := flagSet.String("cert-file", "", "Certificate file")
	keyFile := flagSet.String("key-file", "", "Key file")

	flagSet.Parse(args[2:])

	server.Listen(server.WasmArguments{
		FilePath:        wasmFilePath,
		FunctionName:    wasmFunctionName,
		HTTPPort:        *httpPort,
		Input:           *input,
		LogLevel:        *logLevel,
		AllowHosts:      *allowHosts,
		AllowPaths:      *allowPaths,
		Config:          *config,
		Wasi:            *wasi,
		URL:             *wasmURL,
		AuthHeaderName:  *authHeaderName,
		AuthHeaderValue: *authHeaderValue,
		CertFile:        *certFile,
		KeyFile:         *keyFile,
	})
}
