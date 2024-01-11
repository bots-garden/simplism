package spawn

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	httpHelper "simplism/helpers/http"
	yamlHelper "simplism/helpers/yaml"
	"simplism/server/discovery"
	"simplism/server/router"
	simplismTypes "simplism/types"
	"strconv"

	processesHelper "simplism/helpers/processes"
	stringHelper "simplism/helpers/stringHelper"
)

// This map will store the spawned processes
// It will be used to generate a recovery yam file
var spawnedProcesses = map[string]simplismTypes.WasmArguments{}

var NotifySpawnServiceForRecovery func(formerProcessesArguments map[string]simplismTypes.WasmArguments)

// GetNewHTTPPort returns a unique http port
func getNewHTTPPort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	httpPort := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	listener.Close()
	return httpPort
}

func restartWasmProcess(processArgs simplismTypes.WasmArguments) {
	go func() {
		processesHelper.SpawnSimplismProcess(processArgs)
	}()
}

// spawnHandler returns an http.HandlerFunc that handles requests to spawn a new instance.
//
// It takes wasmArgs simplismTypes.WasmArguments as a parameter.
// It returns an http.HandlerFunc.
func Handler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

	notifyForRecovery := func(formerProcessesArguments map[string]simplismTypes.WasmArguments) {
		fmt.Println("⏳ [recovery] restarting the previous processes")
		//fmt.Println(formerProcessesArguments)
		// Loop through the map
		for _, processArgs := range formerProcessesArguments {
			fmt.Println("🏁 starting:", processArgs.ServiceName, "...")

			if wasmArgs.HttpPortAuto == true {
				processArgs.HTTPPort = getNewHTTPPort()
			}

			spawnedProcesses[processArgs.HTTPPort] = processArgs
			// save the spawned processes to the recovery file
			yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)

			restartWasmProcess(processArgs)
		}

	}
	NotifySpawnServiceForRecovery = notifyForRecovery

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckSpawnToken(request, wasmArgs)

		switch { // /spawn
		case request.Method == http.MethodPost && authorised == true:

			//TODO: if the service name already exists, change the name by name + http port

			/* Request: Create a new Simplism process:

			curl -X POST \
			http://localhost:8080/spawn \
			-H 'admin-spawn-token:michael-burnham-rocks' \
			-H 'Content-Type: application/json; charset=utf-8' \
			--data-binary @- << EOF
			{
				"wasm-file":"../say-hello/say-hello.wasm",
				"wasm-function":"handle",
				"http-port":"9091",
				"discovery-endpoint":"http://localhost:8080/discovery",
				"admin-discovery-token":"michael-burnham-rocks",
				"information": "✋ I'm listening on port 9091",
				"service-name": "say-hello_9091"
			}
			EOF
			*/
			body := httpHelper.GetBody(request)
			bodyMap := map[string]string{}
			err := json.Unmarshal([]byte(body), &bodyMap)

			if err != nil {
				// send response http code error
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 " + err.Error()))
			} else {
				// ✋ right now, you cannot spawn a new spaner process

				response.WriteHeader(http.StatusOK)
				response.Write([]byte("🚀 spawning process...")) // TODO: should be changed
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
				yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)

				// TODO: send the status, only if the process is started (if it's possible)
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
			/* Request: Kill a Simplism process:
			curl -X DELETE \
			http://localhost:8080/spawn?simplismid=42 \
			-H 'admin-spawn-token:michael-burnham-rocks'

			or
			curl -X DELETE \
			http://localhost:8080/spawn?simplismid=42&simplismid=34&simplismid=78 \
			-H 'admin-spawn-token:michael-burnham-rocks'

			*/
			query := request.URL.Query()

			simplismIdList, present := query["simplismid"]
			if !present || len(simplismIdList) == 0 {
				response.WriteHeader(http.StatusNotFound)
				response.Write([]byte("simplismid not present"))
			} else {
				//pid, err := strconv.Atoi(s)
				for _, simplismId := range simplismIdList {
					pid, err := strconv.Atoi(simplismId)
					if err != nil {
						// do nothing
					} else {
						// kill the process
						errKill := processesHelper.KillSimplismProcess(pid)
						if errKill != nil {
							fmt.Println("😡 handler-spawn/KillSimplismProcess", errKill)
						} else {

							foundProcess, errKillNotification := discovery.NotifyDiscoveryServiceOfKillingProcess(pid)

							// Update the recovery file (remove the entry for the killed process)
							delete(spawnedProcesses, foundProcess.HTTPPort)
							yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)

							// Change the handler
							router.GetRouter().HandleFunc("/service/"+foundProcess.ServiceName, func(response http.ResponseWriter, request *http.Request) {
								response.WriteHeader(http.StatusNotFound)
								response.Write([]byte("(Not found) Simplism processe killed"))
							})

							fmt.Println("🙂 Process killed successfully:", foundProcess.ServiceName)

							if errKillNotification != nil {
								fmt.Println("😡 handler-spawn/NotifyDiscoveryServiceOfKillingProcess", errKillNotification)
							} else {
								fmt.Println("🙂 Notification for process killed sent for db update")
							}
						}
					}
					//? Question: kill only one process (? 🤔)
				}

				response.WriteHeader(http.StatusOK)
				response.Write([]byte("Simplism processe(s) killed"))
			}

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}
	}
}
