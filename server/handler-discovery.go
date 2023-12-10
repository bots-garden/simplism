package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpHelper "simplism/helpers/http"
	jsonHelper "simplism/helpers/json"
	simplismTypes "simplism/types"
)

var NotifyDiscoveryServiceOfKillingProcess func(pid int) error

// discoveryHandler handles the /discovery endpoint in the API.
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
func discoveryHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {
	fmt.Println("🔎 discovery mode activated: /discovery  (", wasmArgs.HTTPPort, ")")

	db, _ := initializeDB(wasmArgs)
	// TODO: look at old records and delete old ones

	// This function is called by the spawn handler (DELETE method), see handle-spawn.go
	notifyForKill := func(pid int) error {
		simplismProcess := getSimplismProcessByPiD(db, pid)

		// test simplismProcess.StopTime
		if simplismProcess.StopTime.IsZero() {
		    fmt.Println("⏳ Stop time is not set")
			simplismProcess.StopTime = time.Now()

			err := saveSimplismProcessToDB(db, simplismProcess)
			if err != nil {
				fmt.Println("😡 When updating bucket with the Stop Time", err)
				
			} else {
				fmt.Println("🙂 Bucket updated with the Stop Time")
			}
			return err

		} else {
		    fmt.Println("⏳ Stop time:", simplismProcess.StopTime)
			fmt.Println("✋ This process is already killed")
		}
		
		return nil

	}
	NotifyDiscoveryServiceOfKillingProcess = notifyForKill

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckDiscoveryToken(request, wasmArgs)

		switch {
		// triggered when a simplism process contacts the discovery endpoint
		case request.Method == http.MethodPost && authorised == true:

			body := httpHelper.GetBody(request) // process information from simplism POST request

			// store the process information in the database
			simplismProcess, _ := jsonHelper.GetSimplismProcesseFromJSONBytes(body)
			err := saveSimplismProcessToDB(db, simplismProcess)

			if err != nil {
				fmt.Println("😡 When updating bucket", err)
				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.WriteHeader(http.StatusOK)
			}

		case request.Method == http.MethodGet && authorised == true:

			// get the list of the services that are running
			processes := getSimplismProcessesListFromDB(db)
			jsonString, err := json.Marshal(processes)

			if err != nil {
				fmt.Println("😡 When marshalling", err)
				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.WriteHeader(http.StatusOK)
				response.Write(jsonString)
			}

		case request.Method == http.MethodPut && authorised == true:
			// TODO update the Information field of the service
			// if the token is propagated, the service will be able to PUT information

		// to kill a service, see the admin handler

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			//fmt.Println("😡 You're not authorized")
			//fmt.Fprintln(response, "😡 You're not authorized")
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
			//fmt.Fprintln(response, "😡 Method not allowed")
		}

	}

}
