package discovery

import (
	"fmt"
	"net/http"

	httpHelper "simplism/helpers/http"
	data "simplism/server/data"
	simplismTypes "simplism/types"
)

var NotifyProcessKilled func(pid int) (simplismTypes.SimplismProcess, error)
var NotifyProcessAsleep func(pid int) (simplismTypes.SimplismProcess, error)
var NotifyProcesseInformation func(serviceName string) (simplismTypes.SimplismProcess, error)

func Handler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {
	fmt.Println("🔎 discovery mode activated: /discovery  (", wasmArgs.HTTPPort, ")")

	db, _ := data.InitializeProcessesDB(wasmArgs)

	// Define listeners
	NotifyProcessKilled = processKilledListener(db)
	NotifyProcessAsleep = processAsleepListener(db)
	NotifyProcesseInformation = processInformationListener(db)

	// return the discovery handler
	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckDiscoveryToken(request, wasmArgs)

		switch {

		// ----------------------------
		// Registration
		// ----------------------------
		case request.Method == http.MethodPost && authorised == true:
			// triggered when a simplism process contacts the discovery endpoint
			// a simplism process has been found
			// and try to register to the discovery service
			// see go-routine-simplism-process.go

			simplismProcess, err := registerProcess(request, db)

			if err != nil {
				sendInternalServerErrorResponse(response, err)
			} else {
				createServiceHandler(simplismProcess)
				response.WriteHeader(http.StatusOK)
			}

		// ----------------------------
		// Get the list of processes
		// ----------------------------
		case request.Method == http.MethodGet && authorised == true:

			switch {
			case httpHelper.IsJsonContent(request):

				jsonData, err := getJSONProcesses(db)
				sendJSonResponse(response, jsonData, err)

			case httpHelper.IsTextContent(request):

				data, err := getTableProcesses(db)
				sendTableResponse(response, data, err)
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
