package spawn

import (
	"encoding/json"
	"net/http"
	httpHelper "simplism/helpers/http"
	processesHelper "simplism/helpers/processes"
	stringHelper "simplism/helpers/stringHelper"
	yamlHelper "simplism/helpers/yaml"
	simplismTypes "simplism/types"
)

func subHandlerSpawnProcess(request *http.Request, response http.ResponseWriter, wasmArgs simplismTypes.WasmArguments) {

	body := httpHelper.GetBody(request)
	bodyMap := map[string]string{}
	err := json.Unmarshal([]byte(body), &bodyMap)

	if err != nil {
		// send response http code error
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("😡 " + err.Error()))
	} else {
		// ✋ right now, you cannot spawn a new spaner process

		// ! Start the new process here
		wasmArgsFromJsonPayload := simplismTypes.WasmArguments{}

		wasmArgsFromJsonPayload.FilePath = bodyMap["wasm-file"]
		wasmArgsFromJsonPayload.FunctionName = bodyMap["wasm-function"]
		wasmArgsFromJsonPayload.URL = bodyMap["wasm-url"]
		wasmArgsFromJsonPayload.WasmURLAuthHeader = bodyMap["wasm-url-auth-header"]

		// Automatically assign an HTTP port number to the new process
		if wasmArgs.HttpPortAuto == true {
			wasmArgsFromJsonPayload.HTTPPort = getNewHTTPPort()
		} else {
			wasmArgsFromJsonPayload.HTTPPort = bodyMap["http-port"]
		}

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

		wasmArgsFromJsonPayload.Information = bodyMap["information"]
		wasmArgsFromJsonPayload.ServiceName = bodyMap["service-name"]

		wasmArgsFromJsonPayload.StoreMode = stringHelper.GetTheBooleanValueOf(bodyMap["store-mode"])
		wasmArgsFromJsonPayload.StorePath = bodyMap["store-path"]
		wasmArgsFromJsonPayload.AdminStoreToken = bodyMap["admin-store-token"]

		wasmArgsFromJsonPayload.RegistryMode = stringHelper.GetTheBooleanValueOf(bodyMap["registry-mode"])
		wasmArgsFromJsonPayload.RegistryPath = bodyMap["registry-path"]
		wasmArgsFromJsonPayload.AdminRegistryToken = bodyMap["admin-registry-token"]
		wasmArgsFromJsonPayload.PrivateRegistryToken = bodyMap["private-registry-token"]

		// for debugging
		//fmt.Println("🤓", wasmArgsFromJsonPayload.Information, wasmArgsFromJsonPayload.ServiceName)

		spawnedProcesses[wasmArgsFromJsonPayload.HTTPPort] = wasmArgsFromJsonPayload

		// TODO: handle the error(s) here
		// save the spawned processes to the recovery file
		//yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)
		yamlHelper.WriteYamlFile(GetRecoveryPath(), spawnedProcesses)

		// TODO: send the status, only if the process is started (if it's possible)
		go func() {
			processesHelper.SpawnSimplismProcess(wasmArgsFromJsonPayload)
		}()

		response.WriteHeader(http.StatusOK)
		response.Write([]byte("🚀 process spawned: " + wasmArgsFromJsonPayload.ServiceName))
		// TODO: return json?
	}

}
