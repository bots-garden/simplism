package discovery

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	httpHelper "simplism/helpers/http"
	jsonHelper "simplism/helpers/json"

	//"simplism/server"
	simplismTypes "simplism/types"

	data "simplism/server/data"
	"simplism/server/router"
	//"github.com/go-chi/chi/v5"
)

var NotifyProcessKilled func(pid int) (simplismTypes.SimplismProcess, error)
var NotifyGetProcesseInformation func(serviceName string) (simplismTypes.SimplismProcess, error)

var wasmFunctionHandlerList = map[string]int{}

// This variable is initialized at start if the server is in discovery mode
// It stores the processes that was previously running
// Then we can restart them in spawn mode
// TODO: use this to restart after a crash?
//var formerProcesses  map[string]simplismTypes.SimplismProcess

// DiscoveryHandler handles the /discovery endpoint in the API.
//
// It takes a WasmArguments object as a parameter and returns an http.HandlerFunc.
// The WasmArguments object contains information about the HTTP port.
// The returned http.HandlerFunc handles incoming HTTP requests to the /discovery endpoint.
// It checks if the request is authorized and if it is a POST request.
// If authorized and a POST request, it processes the information from the request body,
// creates a SimpleProcess struct instance from the JSON body, and stores the process information in the database.
// If there is an error while saving the process information, it returns a 500 Internal Server Error response.
// If the request is not authorized, it returns a 401 Unauthorized response.
// If the request method is not allowed, it returns a 405 Method Not Allowed response.
// This function is a work in progress and handles GET, DELETE, and PUT requests.
func Handler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {
	fmt.Println("🔎 discovery mode activated: /discovery  (", wasmArgs.HTTPPort, ")")

	db, _ := data.InitializeProcessesDB(wasmArgs)
	// TODO: look at old records and delete old ones

	//formerProcesses = getSimplismProcessesListFromDB(db)

	// This function is called by the spawn handler (DELETE method), see handle-spawn.go
	notifyProcessKilled := func(pid int) (simplismTypes.SimplismProcess, error) {
		simplismProcess := data.GetSimplismProcessByPiD(db, pid)

		if simplismProcess.PID == 0 {
			return simplismTypes.SimplismProcess{}, errors.New("😡 Process not found")
		} else {
			// test simplismProcess.StopTime
			if simplismProcess.StopTime.IsZero() {
				fmt.Println("⏳ Stop time is not set")
				simplismProcess.StopTime = time.Now()

				err := data.SaveSimplismProcessToDB(db, simplismProcess)
				if err != nil {
					fmt.Println("😡 When updating bucket with the Stop Time", err)
					return simplismTypes.SimplismProcess{}, err

				} else {
					fmt.Println("🙂 Bucket updated with the Stop Time")
				}

			} else {
				fmt.Println("⏳ Stop time:", simplismProcess.StopTime)
				fmt.Println("✋ This process is already killed")
			}

			return simplismProcess, nil
		}

	}
	NotifyProcessKilled = notifyProcessKilled

	notifyGetProcessInformation := func(serviceName string) (simplismTypes.SimplismProcess, error) {

		simplismProcess := data.GetSimplismProcessByName(db, serviceName)
		
		if simplismProcess.PID == 0 {
			return simplismTypes.SimplismProcess{}, errors.New("😡 Process not found")
		}
		return simplismProcess, nil
	}
	NotifyGetProcesseInformation = notifyGetProcessInformation

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckDiscoveryToken(request, wasmArgs)

		switch {
		// triggered when a simplism process contacts the discovery endpoint
		case request.Method == http.MethodPost && authorised == true:

			body := httpHelper.GetBody(request) // process information from simplism POST request

			// store the process information in the database
			simplismProcess, _ := jsonHelper.GetSimplismProcesseFromJSONBytes(body)
			err := data.SaveSimplismProcessToDB(db, simplismProcess)

			//simplismProcess.ServiceName

			if err != nil {
				fmt.Println("😡 When updating bucket", err)
				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.WriteHeader(http.StatusOK)

				/* Call a function from the discovery service
				   ------------------------------------------
					if there is a new simplism function process contact
					- create a new handler to handle the requests (kind of reverse proxy)
					- only if the handler doesn't exist

					if the process service name is "hello" and listening on port 9090
					if the process spawaner is listening on port 8080

					when you call http://localhost:8080/function/hello
					a request will be sent to http://localhost:9090/function/hello

				*/

				if wasmFunctionHandlerList[simplismProcess.ServiceName] == 0 {
					wasmFunctionHandlerList[simplismProcess.ServiceName] = simplismProcess.PID

					router.GetRouter().HandleFunc("/service/"+simplismProcess.ServiceName, func(response http.ResponseWriter, request *http.Request) {

						host, _, _ := net.SplitHostPort(request.Host)

						// make an HTTP request to the simplismservice
						//! https? handled by the spawner
						client := &http.Client{}
						body := httpHelper.GetBody(request)
						requestToSpawnedProcess, _ := http.NewRequest(request.Method, "http://"+host+":"+simplismProcess.HTTPPort, bytes.NewBuffer(body))
						requestToSpawnedProcess.Header = request.Header

						// Send the request
						responseFromSpawnedProcess, err := client.Do(requestToSpawnedProcess)
						if err != nil {
							fmt.Println("😡 When making the HTTP request", err)
						}
						defer responseFromSpawnedProcess.Body.Close()
						// Read the response body
						responseBodyFromSpawnedProcess, err := io.ReadAll(responseFromSpawnedProcess.Body)
						if err != nil {
							fmt.Println("😡 Error reading response body:", err)
							return
						}

						response.WriteHeader(responseFromSpawnedProcess.StatusCode)
						response.Write(responseBodyFromSpawnedProcess)

					})
				}

			}

		case request.Method == http.MethodGet && authorised == true:

			switch {
			case httpHelper.IsJsonContent(request):

				jsonData, err := getJSONProcesses(db)
				sendJSonResponse(response, jsonData, err)

			case httpHelper.IsTextContent(request):

				data, err := getTableProcesses(db)
				sendTableResponse(response, data, err)
			}

		//case request.Method == http.MethodPut && authorised == true:
		// TODO update the Information field of the service

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}

	}

}
