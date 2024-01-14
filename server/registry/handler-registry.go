package registry

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	httpHelper "simplism/helpers/http"
	simplismTypes "simplism/types"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Handler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

	return func(response http.ResponseWriter, request *http.Request) {

		adminRegistryAuthorised := httpHelper.CheckAdminRegistryToken(request, wasmArgs)
		privateRegistryAuthorised := httpHelper.CheckPrivateRegistryToken(request, wasmArgs)

		switch { // /registry
		// upload: /registry/push
		case request.Method == http.MethodPost && adminRegistryAuthorised == true:
			// Maximum upload of 10 MB files
			request.ParseMultipartForm(10 << 20)

			// Get handler for filename, size and headers
			file, handler, err := request.FormFile("file")
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error Retrieving the File"))
				return
			}

			defer file.Close()

			fmt.Printf("📝 uploaded file: %+v\n", wasmArgs.RegistryPath+"/"+handler.Filename)
			fmt.Printf("📐 file size: %+v\n", handler.Size)
			fmt.Printf("📙 MIME header: %+v\n", handler.Header)

			// Create file
			dst, err := os.Create(wasmArgs.RegistryPath + "/" + handler.Filename)
			defer dst.Close()
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error creating wasm file"))
				return
			}

			// Copy the uploaded file to the created file on the filesystem
			if _, err := io.Copy(dst, file); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error copying wasm file"))
				return
			}

			response.WriteHeader(http.StatusOK)
			response.Write([]byte("🎉 "+ handler.Filename +" successfully uploaded!\n"))

		// download: /registry/pull/{wasmfilename}
		case request.Method == http.MethodGet &&
			strings.HasPrefix(request.RequestURI, "/registry/pull/") &&
			strings.HasSuffix(request.RequestURI, ".wasm") &&
			privateRegistryAuthorised == true:

			filename := chi.URLParam(request, "wasmfilename")

			f, err := os.Open(wasmArgs.RegistryPath + "/" + filename)
			if err != nil {
				fmt.Println("😡 Error reading wasm file")
				response.WriteHeader(http.StatusNotFound)
				response.Write([]byte("😡 Error reading wasm file"))
				return
			}

			wasmFile, err := io.ReadAll(f)
			if err != nil {
				panic(err)
			}

			setHeaders(response, filename, strconv.Itoa(len(wasmFile)))
			//response.WriteHeader(http.StatusOK)
			n, err := response.Write(wasmFile)

			if err != nil {
				//fmt.Println("😡 Error writing wasm file")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error writing wasm file"))
				return
			}
			log.Printf("📝 wasm file written : %d", n)

		// /registry/discover
		case request.Method == http.MethodGet &&
			strings.HasPrefix(request.RequestURI, "/registry/discover") &&
			privateRegistryAuthorised == true:

			switch {
			case httpHelper.IsJsonContent(request):
				jsonData, err := getJSONListOfFiles(wasmArgs.RegistryPath)
				sendJSonResponse(response, jsonData, err)

			case httpHelper.IsTextContent(request):
				data, err := getTableListOfFiles(wasmArgs.RegistryPath)
				sendTableResponse(response, data, err)
			}


		// /registry/remove/{wasmfilename}
		case request.Method == http.MethodDelete &&
			strings.HasPrefix(request.RequestURI, "/registry/remove") &&
			strings.HasSuffix(request.RequestURI, ".wasm") &&
			adminRegistryAuthorised == true:

			fileName := chi.URLParam(request, "wasmfilename")

			err := os.Remove(wasmArgs.RegistryPath + "/" + fileName)

			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error when removing wasm file"))
				return
			}

			response.WriteHeader(http.StatusOK)
			//response.Header().Set("Content-Type", "application/json; charset=utf-8")
			response.Write([]byte("🙂 " + chi.URLParam(request, "wasmfilename") + " removed"))

		case adminRegistryAuthorised == false || privateRegistryAuthorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}

	}
}
