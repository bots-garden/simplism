package registry

import (
	"fmt"
	"net/http"
)

// sendJSonResponse sends a JSON response to the client.
//
// It takes three parameters:
// - response: the http.ResponseWriter used to send the response.
// - jsonData: a byte slice containing the JSON data to be sent.
// - err: an error that occurred when reading files.
//
// The function does the following:
// - If an error occurred, it prints an error message and sets the response status to 500 (Internal Server Error).
// - Otherwise, it sets the response status to 200 (OK), sets the content type to "application/json; charset=utf-8",
//   and writes the JSON data to the response body.
func sendJSonResponse(response http.ResponseWriter, jsonData []byte, err error) {
	if err != nil {
		fmt.Println("😡 When reading files", err)
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
// - data: a 2D slice of strings representing the tabular data.
// - err: an error, if any, that occurred while getting the tabular data.
//
// This function does not return anything.
func sendTableResponse(response http.ResponseWriter, data [][]string ,err error) {
    if err != nil {
        fmt.Println("😡 When getting tabular data", err)
        response.WriteHeader(http.StatusInternalServerError)
    } else {
        table := getListOfFilesTableWriter(response, data)
        response.WriteHeader(http.StatusOK)
        response.Header().Set("Content-Type", "plain/text; charset=utf-8")
        table.Render()
    }
}
