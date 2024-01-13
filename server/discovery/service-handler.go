package discovery

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	httpHelper "simplism/helpers/http"
	"simplism/server/router"
	simplismTypes "simplism/types"
)

// list of the handlers of the services
var wasmFunctionHandlerList = map[string]int{}

/* Call a function from the discovery service
   ------------------------------------------
	if there is a new simplism function process registration
	- create a new handler to handle the requests (kind of reverse proxy)
	- only if the handler doesn't exist

	if the process service name is "hello" and listening on port 9090
	if the process spawaner is listening on port 8080

	when you call http://localhost:8080/service/hello
	a request will be sent to http://localhost:9090/service/hello

*/

func createServiceHandler(simplismProcess simplismTypes.SimplismProcess) {

	// Create the handler for the service only if the handler does not exist
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
