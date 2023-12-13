package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	httpHelper "simplism/helpers/http"
	wasmHelper "simplism/helpers/wasm"
	simplismTypes "simplism/types"
)

/*
	This handler is responsible for:
	- handling HTTP requests and,
	- calling the WebAssembly function.

*/

// mainHandler is the main handler function for HTTP requests.
//
// It takes in the WasmArguments struct as a parameter and returns
// an http.HandlerFunc. Within the returned handler function, it
// extracts the body from the HTTP request, creates a mainFunctionArgument
// using the request's header, body, method, and URI, and calls the
// wasmhelper.CallWasmFunction with the wasmArgs.FunctionName and
// mainFunctionArgument. The result is unmarshaled into a ReturnValue
// struct, and if there is an error during unmarshaling, an HTTP
// internal server error is returned with the error message. Otherwise,
// the header values from the ReturnValue are set in the response
// header. If there is an error during the wasm function call, an
// HTTP internal server error is returned with the error message.
// Otherwise, the response status code is set to the ReturnValue.Code
// and the response body is set to the ReturnValue.Body.
//
// Parameters:
// - wasmArgs: A struct containing the arguments for the Wasm function.
//
// Return Type:
// - http.HandlerFunc: The handler function for HTTP requests.
func mainHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

	return func(response http.ResponseWriter, request *http.Request) {

		var (
			result []byte
			err    error
		)

		body := httpHelper.GetBody(request) // is the body the same with fiber ?

		mainFunctionArgument := simplismTypes.Argument{
			Header: request.Header,
			Body:   string(body),
			Method: request.Method,
			URI:    request.RequestURI,
		}

		//result, err = wasmHelper.CallWasmFunction(wasmFunctionName, []byte(mainFunctionArgument.ToEncodedJSONString()))
		result, err = wasmHelper.CallWasmFunction(wasmArgs.FunctionName, mainFunctionArgument.ToJSONBuffer())

		/* Expected response
		type ReturnValue struct {
			Body   string              `json:"body"`
			Header map[string][]string `json:"header"`
		}
		*/
		returnValue := simplismTypes.ReturnValue{}
		errJSONUnmarshal := json.Unmarshal(result, &returnValue)

		if errJSONUnmarshal != nil {
			// send response http code error
			response.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(response, errJSONUnmarshal.Error())
		} else {
			for key, value := range returnValue.Header {
				response.Header().Set(key, value[0])
			}

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(response, err.Error())
			} else {
				// TODO: add default response code if empty?
				response.WriteHeader(returnValue.Code)
				fmt.Fprintln(response, string(returnValue.Body))
			}
		}

	}

}
