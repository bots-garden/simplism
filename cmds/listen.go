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

	// these environment variables are forwarded to the wasm plug-in (not the other variables)
	envVars := flagSet.String("env", `[]`, "Environment variables to forward to the wasm plug-in")

	config := flagSet.String("config", "{}", "Configuration data (json string)")

	wasi := flagSet.Bool("wasi", true, "")

	wasmURL := flagSet.String("wasm-url", "", "Url to download the wasm file")

	wasmURLAuthHeader := flagSet.String("wasm-url-auth-header", "", "Authentication header ex: `PRIVATE-TOKEN=IlovePandas`")
	// or you can use this environment variable: WASM_URL_AUTH_HEADER="PRIVATE-TOKEN=IlovePandas"

	//authHeaderName := flagSet.String("auth-header-name", "", "Authentication header name, ex: PRIVATE-TOKEN")
	//authHeaderValue := flagSet.String("auth-header-value", "", "Value of the authentication header, ex: IlovePandas")

	certFile := flagSet.String("cert-file", "", "Certificate file")
	keyFile := flagSet.String("key-file", "", "Key file")

	/* --- admin tokens --- */

	// admin-reload-token or environment variable: ADMIN_RELOAD_TOKEN
	adminReloadToken := flagSet.String("admin-reload-token", "", "Admin reload token")

	discovery := flagSet.Bool("discovery", false, "")
    discoveryEndpoint := flagSet.String("discovery-endpoint", "", "Discovery endpoint")

	flagSet.Parse(args[2:])

	server.Listen(server.WasmArguments{
		FilePath:          wasmFilePath,
		FunctionName:      wasmFunctionName,
		HTTPPort:          *httpPort,
		Input:             *input,
		LogLevel:          *logLevel,
		AllowHosts:        *allowHosts,
		AllowPaths:        *allowPaths,
		EnvVars:           *envVars,
		Config:            *config,
		Wasi:              *wasi,
		URL:               *wasmURL,
		WasmURLAuthHeader: *wasmURLAuthHeader,
		//AuthHeaderName:  *authHeaderName,
		//AuthHeaderValue: *authHeaderValue,
		CertFile:         *certFile,
		KeyFile:          *keyFile,
		AdminReloadToken: *adminReloadToken,
		Discovery:        *discovery,
		DiscoveryEndpoint: *discoveryEndpoint,
	}, "") // no config key
}
