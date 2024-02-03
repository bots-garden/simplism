package spawn

import (
	"fmt"
	"net"
	"net/http"
	httpHelper "simplism/helpers/http"
	yamlHelper "simplism/helpers/yaml"
	simplismTypes "simplism/types"
	"strconv"
	"strings"

	processesHelper "simplism/helpers/processes"
)

// This map will store the spawned processes
// It will be used to generate a recovery yam file
var spawnedProcesses = map[string]simplismTypes.WasmArguments{}

// Default
var recoveryPath string = "recovery.yaml"

func GetRecoveryPath() string {
	return recoveryPath
}
func SetRecoveryPath(path string) {
	recoveryPath = path
}

var NotifyStartRecovery func(formerProcessesArguments map[string]simplismTypes.WasmArguments)

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

	notifyStartRecovery := func(formerProcessesArguments map[string]simplismTypes.WasmArguments) {
		fmt.Println("⏳ [recovery] restarting the previous processes")
		// Loop through the map
		for _, processArgs := range formerProcessesArguments {
			fmt.Println("🏁 starting:", processArgs.ServiceName, "...")

			if wasmArgs.HttpPortAuto == true {
				processArgs.HTTPPort = getNewHTTPPort()
			}

			spawnedProcesses[processArgs.HTTPPort] = processArgs
			// save the spawned processes to the recovery file

			//yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)
			yamlHelper.WriteYamlFile(GetRecoveryPath(), spawnedProcesses)

			restartWasmProcess(processArgs)
		}

	}
	NotifyStartRecovery = notifyStartRecovery

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
			subHandlerSpawnProcess(request, response, wasmArgs)

		/*
			case request.Method == http.MethodGet && authorised == true:
				response.WriteHeader(http.StatusOK)
				response.Write([]byte("👋 Hello [GET]"))

			case request.Method == http.MethodPut && authorised == true:
				response.WriteHeader(http.StatusOK)
				response.Write([]byte("👋 Hello [PUT]"))
		*/

		//--------------------------------------------------------------
		// Kill or asleep a Simplism process by pid or name:
		//--------------------------------------------------------------
		case request.Method == http.MethodDelete && authorised == true:

			switch {
			/*
				curl -X DELETE \
				http://localhost:8080/spawn/kill/name/hello \
				-H 'admin-spawn-token:michael-burnham-rocks'
			*/
			case strings.HasPrefix(request.RequestURI, "/spawn/kill/name/"):
				subHandlerKillByName(request, response)

			/*
				curl -X DELETE \
				http://localhost:8080/spawn/kill/pid/42 \
				-H 'admin-spawn-token:michael-burnham-rocks'
			*/
			case strings.HasPrefix(request.RequestURI, "/spawn/kill/pid/"):
				subHandlerKillByPid(request, response)

			/*
				curl -X DELETE \
				http://localhost:8080/spawn/fall-asleep/name/hello \
				-H 'admin-spawn-token:michael-burnham-rocks'
			*/
			case strings.HasPrefix(request.RequestURI, "/spawn/fall-asleep/name/"):
				subHandlerFallAsleepByName(request, response)

			/*
				curl -X DELETE \
				http://localhost:8080/spawn/fall-asleep/pid/42 \
				-H 'admin-spawn-token:michael-burnham-rocks'
			*/
			case strings.HasPrefix(request.RequestURI, "/spawn/fall-asleep/pid/"):
				subHandlerFallAsleepByPid(request, response)

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
