package server

import (
	"net/http"
	simplismTypes "simplism/types"
	httpHelper "simplism/helpers/http"

)

func storeHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {
	
	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckStoreToken(request, wasmArgs)

		switch { // /store
		case request.Method == http.MethodPost && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("📦 Hello [POST]"))

		case request.Method == http.MethodGet && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("📦 Hello [GET]"))

		case request.Method == http.MethodPut && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("📦 Hello [PUT]"))

		case request.Method == http.MethodDelete && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("📦 Hello [DELETE]"))

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}
	}
}
