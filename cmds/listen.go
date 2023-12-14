package cmds

import (
	"flag"
	stringHelper "simplism/helpers/stringHelper"
	"simplism/server"
	simplismTypes "simplism/types"
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

	wasi := flagSet.String("wasi", "true", "")

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

	serviceDiscovery := flagSet.String("service-discovery", "false", "")

	discoveryEndpoint := flagSet.String("discovery-endpoint", "", "Discovery endpoint")

	adminDiscoveryToken := flagSet.String("admin-discovery-token", "", "Admin discovery token")

	spawnMode := flagSet.String("spawn-mode", "false", "")
	adminSpawnToken := flagSet.String("admin-spawn-token", "", "Admin spawn token")

	information := flagSet.String("information", "", "Information about the simplism service (useful for the discovery mode)")
	serviceName := flagSet.String("service-name", "", "Simplism service name (useful for the discovery mode)")

	storeMode := flagSet.String("store-mode", "false", "")
	adminStoreToken := flagSet.String("admin-store-token", "", "Admin store token")


	flagSet.Parse(args[2:])

	server.Listen(simplismTypes.WasmArguments{
		FilePath:          wasmFilePath,
		FunctionName:      wasmFunctionName,
		HTTPPort:          *httpPort,
		Input:             *input,
		LogLevel:          *logLevel,
		AllowHosts:        *allowHosts,
		AllowPaths:        *allowPaths,
		EnvVars:           *envVars,
		Config:            *config,
		Wasi:              stringHelper.GetTheBooleanValueOf(*wasi),
		URL:               *wasmURL,
		WasmURLAuthHeader: *wasmURLAuthHeader,
		//AuthHeaderName:  *authHeaderName,
		//AuthHeaderValue: *authHeaderValue,
		CertFile:            *certFile,
		KeyFile:             *keyFile,
		AdminReloadToken:    *adminReloadToken,
		ServiceDiscovery:    stringHelper.GetTheBooleanValueOf(*serviceDiscovery),
		DiscoveryEndpoint:   *discoveryEndpoint,
		AdminDiscoveryToken: *adminDiscoveryToken,
		SpawnMode:           stringHelper.GetTheBooleanValueOf(*spawnMode),
		AdminSpawnToken:     *adminSpawnToken,
		Information:         *information,
		ServiceName:         *serviceName,
		StoreMode:           stringHelper.GetTheBooleanValueOf(*storeMode),
		AdminStoreToken:     *adminStoreToken,
	}, "") // no config key
}
