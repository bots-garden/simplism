package discovery

import (
	"fmt"
	"net/http"
)

// sendJSonResponse sends a JSON response to the client.
//
// The function takes in three parameters:
// - response: an http.ResponseWriter object that is used to write the response back to the client.
// - jsonData: a byte slice containing the JSON data to be sent.
// - err: an error object that represents any error that occurred during the marshalling of the JSON data.
//
// The function returns no values.
func sendJSonResponse(response http.ResponseWriter, jsonData []byte ,err error) {
	if err != nil {
		fmt.Println("😡 When marshalling", err)
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Write(jsonData)
	}
}

// sendTableResponse sends a table response to the client.
//
// It takes three parameters:
// - response: an http.ResponseWriter to write the response to.
// - data: a 2D string slice representing the tabular data.
// - err: an error indicating any errors that occurred while getting the tabular data.
//
// It does not return anything.
func sendTableResponse(response http.ResponseWriter, data [][]string ,err error) {
    if err != nil {
        fmt.Println("😡 When getting tabular data", err)
        response.WriteHeader(http.StatusInternalServerError)
    } else {
        table := getProcessesTableWriter(response, data)
        response.WriteHeader(http.StatusOK)
        response.Header().Set("Content-Type", "plain/text; charset=utf-8")
        table.Render()
    }
}

func sendInternalServerErrorResponse(response http.ResponseWriter, err error) {
	fmt.Println("😡 When updating bucket", err)
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("😡 When updating bucket"))
}