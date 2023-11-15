package main

// create a simple http server
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	check "simplism/checkers"
	wasmHelper "simplism/extism-runtime"
	functionTypes "simplism/function-types"
	httpHelper "simplism/http-helper"
)

func main() {

	args := os.Args[1:]

	// Exit if not enough args
	check.IfThereAreEnoughArgs(args)

	wasmFilePath := args[0]
	wasmFunctionName := args[1]
	httpPort := args[2]

	ctx := context.Background()

	config, manifest := wasmHelper.GetConfigAndManifest(wasmFilePath)

	wasmHelper.GeneratePluginsPool(ctx, config, manifest)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		var (
			result []byte
			err    error
		)

		body := httpHelper.GetBody(request) // is the body the same with fiber ?

		mainFunctionArgument := functionTypes.Argument{
			Header: request.Header,
			Body:   string(body),
			Method: request.Method,
			URI:    request.RequestURI,
		}

		//result, err = wasmHelper.CallWasmFunction(wasmFunctionName, []byte(mainFunctionArgument.ToEncodedJSONString()))
		result, err = wasmHelper.CallWasmFunction(wasmFunctionName, mainFunctionArgument.ToJSONBuffer())

		/* Expected response
		type ReturnValue struct {
			Body   string              `json:"body"`
			Header map[string][]string `json:"header"`
		}
		*/
		returnValue := functionTypes.ReturnValue{}
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

	})

	go func() {
		// certificate - https ?
		fmt.Println("🌍 http server is listening on:", httpPort)
		log.Fatal(http.ListenAndServe(":"+httpPort, nil))
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
}
