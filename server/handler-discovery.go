package server

import (
	"fmt"
	"net/http"
	"os"
	"simplism/httphelper"
	"simplism/jsonhelper"

	simplismTypes "simplism/types"
)

// checkDiscoveryToken checks if the provided request is authorized using the admin-discovery-token.
//
// Parameters:
// - request: the HTTP request to check.
// - wasmArgs: the Wasm arguments.
//
// Returns:
// - bool: true if the request is authorized, false otherwise.
func checkDiscoveryToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false
	// read the header admin-discovery-token
	adminDiscoveryToken := request.Header.Get("admin-discovery-token")

	if wasmArgs.AdminDiscoveryToken != "" {
		// token is awaited
		if wasmArgs.AdminDiscoveryToken == adminDiscoveryToken {
			authorised = true
		} else {
			authorised = false
		}
	} else {
		// check if the env variable ADMIN_DISCOVERY_TOKEN is set
		envAdminDiscoveryToken := os.Getenv("ADMIN_DISCOVERY_TOKEN")
		if envAdminDiscoveryToken != "" {
			// token is awaited
			if envAdminDiscoveryToken == adminDiscoveryToken {
				authorised = true
			} else {
				authorised = false
			}
		} else {
			authorised = true
		}

	}

	return authorised
}

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

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := checkDiscoveryToken(request, wasmArgs)
		// Test if it is a POST request
		if request.Method == http.MethodPost && authorised == true {

			body := httphelper.GetBody(request) // process information from simplism POST request

			// create SimpleProcess struct instance from JSON Body
			simplismProcess, _ := jsonhelper.GetSimplismProcesseFromJSONBytes(body)
			// store the process information in the database
			err := saveSimplismProcessToDB(db, simplismProcess)

			// TODO: look at old records and delete old ones
			if err != nil {
				fmt.Println("😡 When updating bucket", err)
				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.WriteHeader(http.StatusOK)
			}

		} else {
			if authorised == false {
				response.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(response, "😡 You're not authorized")

			} else {
				response.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintln(response, "😡 Method not allowed")
			}
			// 🚧 This is a Work In Progress
			// If GET request

			// If DELETE request

			// If PUT request
		}

	}

}
