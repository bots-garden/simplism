package reload

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	configHelper "simplism/helpers/config"
	httpHelper "simplism/helpers/http"
	wasmHelper "simplism/helpers/wasm"
	"simplism/server/processes"
	simplismTypes "simplism/types"
)

/*
	This handler is responsible for:
	- reloading the WebAssembly file,
*/

// reloadHandler handles the HTTP request for reloading the wasm plug-in.
//
// It expects the following:
// - A POST request
// - A URL for the wasm file
// - An admin-reload-token value in the request header and the ADMIN_RELOAD_TOKEN environment variable
//
// Parameters:
// - ctx: The context.Context object.
// - wasmArgs: The WasmArguments struct that contains the necessary arguments for reloading the wasm plug-in.
//
// Return type:
// - http.HandlerFunc
func Handler(ctx context.Context, wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// wait for:
		// - POST request
		// - url for the wasm file
		// - admin-reload-token and ADMIN_RELOAD_TOKEN env variable
		authorised := httpHelper.CheckReloadToken(request, wasmArgs)

		// Test if it's a POST request
		if request.Method == http.MethodPost && authorised == true {
			body := httpHelper.GetBody(request)
			// body is a JSON string, extract the url field value of the JSON string
			bodyMap := map[string]string{}
			err := json.Unmarshal([]byte(body), &bodyMap)

			if err != nil {
				// send response http code error
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 " + err.Error()))
				//fmt.Fprintln(response, "😡 "+err.Error())
			} else {

				wasmArgs.URL = bodyMap["wasm-url"]
				wasmArgs.FilePath = bodyMap["wasm-file"]
				wasmArgs.FunctionName = bodyMap["wasm-function"]

				/*
					wasmArgs.WasmURLAuthHeader = bodyMap["wasm-url-auth-header"]
					wasmArgs.LogLevel = bodyMap["log-level"]
					wasmArgs.AllowHosts = bodyMap["allow-hosts"]
					wasmArgs.AllowPaths = bodyMap["allow-paths"]
					wasmArgs.EnvVars = bodyMap["env"]
					wasmArgs.Config = bodyMap["config"]
					wasmArgs.Wasi = stringHelper.GetTheBooleanValueOf(bodyMap["wasi"])
					wasmArgs.Input = bodyMap["input"]
					wasmArgs.CertFile = bodyMap["cert-file"]
					wasmArgs.KeyFile = bodyMap["key-file"]
					wasmArgs.AdminReloadToken = bodyMap["admin-reload-token"]
					wasmArgs.AdminDiscoveryToken = bodyMap["admin-discovery-token"]
					wasmArgs.ServiceDiscovery = stringHelper.GetTheBooleanValueOf(bodyMap["service-discovery"])
					wasmArgs.DiscoveryEndpoint = bodyMap["discovery-endpoint"]

					wasmArgs.Information = bodyMap["information"]
					wasmArgs.ServiceName = bodyMap["service-name"]

					wasmArgs.StoreMode = stringHelper.GetTheBooleanValueOf(bodyMap["store-mode"])
					wasmArgs.StorePath = bodyMap["store-path"]
					wasmArgs.AdminStoreToken = bodyMap["admin-store-token"]

					wasmArgs.RegistryMode = stringHelper.GetTheBooleanValueOf(bodyMap["registry-mode"])
					wasmArgs.RegistryPath = bodyMap["registry-path"]
					wasmArgs.AdminRegistryToken = bodyMap["admin-registry-token"]
					wasmArgs.PrivateRegistryToken = bodyMap["private-registry-token"]
				*/

				fmt.Println("🚀 downloading", wasmArgs.URL, "...")
				err := wasmHelper.DownloadWasmFile(wasmArgs)
				if err != nil {
					fmt.Println(err)
					os.Exit(1) // TODO: do something else
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

				config, manifest := wasmHelper.GetConfigAndManifest(wasmArgs.FilePath, hosts, paths, manifestConfig, level)

				wasmHelper.StartWasmPlugin(ctx, config, manifest)
				fmt.Println("🙂 new wasm plug-in reloaded")

				response.WriteHeader(http.StatusOK)
				//fmt.Fprintln(response, string("🙂 new wasm plug-in reloaded"))
				response.Write([]byte("🙂 new wasm plug-in reloaded"))

				// Update information about the current simplism process
				processes.SetCurrentProcessFilePath(wasmArgs.FilePath)
				processes.SetCurrentProcessFunctionName(wasmArgs.FunctionName)

			}

		} else {
			if authorised == false {
				response.WriteHeader(http.StatusUnauthorized)
				//fmt.Fprintln(response, "😡 You're not authorized")
				response.Write([]byte("😡 You're not authorized"))

			} else {
				response.WriteHeader(http.StatusMethodNotAllowed)
				//fmt.Fprintln(response, "😡 Method not allowed")
				response.Write([]byte("😡 Method not allowed"))
			}
		}

	}
}
