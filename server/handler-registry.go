package server

import (
	"net/http"
	httpHelper "simplism/helpers/http"
	simplismTypes "simplism/types"
)

func registryHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckRegistryToken(request, wasmArgs)

		switch { // /registry
		case request.Method == http.MethodPost && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("🙂 POST"))

		case request.Method == http.MethodGet && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("🙂 GET"))

		case request.Method == http.MethodPut && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("🙂 PUT"))

		case request.Method == http.MethodDelete && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("🙂 DELETE"))
            
		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}

	}
}
