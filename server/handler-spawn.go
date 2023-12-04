package server

import (
	"encoding/json"
	"net/http"
	httpHelper "simplism/helpers/http"
	simplismTypes "simplism/types"

	processesHelper "simplism/helpers/processes"
	stringHelper "simplism/helpers/stringHelper"
)

// spawnHandler returns an http.HandlerFunc that handles requests to spawn a new instance.
//
// It takes wasmArgs simplismTypes.WasmArguments as a parameter.
// It returns an http.HandlerFunc.
func spawnHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {
	// TODO
	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckSpawnToken(request, wasmArgs)

		switch {
		case request.Method == http.MethodPost && authorised == true:

			body := httpHelper.GetBody(request)
			bodyMap := map[string]string{}
			err := json.Unmarshal([]byte(body), &bodyMap)

			if err != nil {
				// send response http code error
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 " + err.Error()))
			} else {
				response.WriteHeader(http.StatusOK)
				response.Write([]byte("🚀 spawning mode, work in progress"))
				// ! Start the new process here
				wasmArgsFromJsonPayload := simplismTypes.WasmArguments{}

				wasmArgsFromJsonPayload.FilePath = bodyMap["wasm-file"]
				wasmArgsFromJsonPayload.FunctionName = bodyMap["wasm-function"]
				wasmArgsFromJsonPayload.URL = bodyMap["wasm-url"]
				wasmArgsFromJsonPayload.WasmURLAuthHeader = bodyMap["wasm-url-auth-header"]
				wasmArgsFromJsonPayload.HTTPPort = bodyMap["http-port"]
				wasmArgsFromJsonPayload.LogLevel = bodyMap["log-level"]
				wasmArgsFromJsonPayload.AllowHosts = bodyMap["allow-hosts"]
				wasmArgsFromJsonPayload.AllowPaths = bodyMap["allow-paths"]
				wasmArgsFromJsonPayload.EnvVars = bodyMap["env"]
				wasmArgsFromJsonPayload.Config = bodyMap["config"]
				wasmArgsFromJsonPayload.Wasi = stringHelper.GetTheBooleanValueOf(bodyMap["wasi"])
				wasmArgsFromJsonPayload.Input = bodyMap["input"]
				wasmArgsFromJsonPayload.CertFile = bodyMap["cert-file"]
				wasmArgsFromJsonPayload.KeyFile = bodyMap["key-file"]
				wasmArgsFromJsonPayload.AdminReloadToken = bodyMap["admin-reload-token"]
				wasmArgsFromJsonPayload.AdminDiscoveryToken = bodyMap["admin-discovery-token"]
				//wasmArgsFromJsonPayload.AdminSpawnToken = bodyMap["admin-spawn-token"]
				wasmArgsFromJsonPayload.ServiceDiscovery = stringHelper.GetTheBooleanValueOf(bodyMap["service-discovery"])
				wasmArgsFromJsonPayload.DiscoveryEndpoint = bodyMap["discovery-endpoint"]

				go func() {
					processesHelper.SpawnSimplismProcess(wasmArgsFromJsonPayload)
				}()

			}

		case request.Method == http.MethodGet && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("👋 Hello [GET]"))

		case request.Method == http.MethodPut && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("👋 Hello [PUT]"))

		case request.Method == http.MethodDelete && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("👋 Hello [DELETE]"))

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))
			//fmt.Fprintln(response, "😡 You're not authorized")

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
			//fmt.Fprintln(response, "😡 Method not allowed")
		}

	}
}
