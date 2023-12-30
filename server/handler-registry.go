package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	httpHelper "simplism/helpers/http"
	simplismTypes "simplism/types"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func setHeaders(response http.ResponseWriter, name, len string) {
	//Represents binary file
	response.Header().Set("Content-Type", "application/octet-stream")
	//Tells client what filename should be used.
	response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	//The length of the data.
	response.Header().Set("Content-Length", len)
	//No cache headers.
	response.Header().Set("Cache-Control", "private")
	//No cache headers.
	response.Header().Set("Pragma", "private")
	//No cache headers.
	response.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func registryHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckRegistryToken(request, wasmArgs)

		switch { // /registry
		// upload
		case request.Method == http.MethodPost && authorised == true:
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
				//http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}

			// Copy the uploaded file to the created file on the filesystem
			if _, err := io.Copy(dst, file); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error copying wasm file"))
				//http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}

			response.WriteHeader(http.StatusOK)
			response.Write([]byte("🎉 Successfully Uploaded File\n"))

		// download
		case request.Method == http.MethodGet && authorised == true:
			//fmt.Printf("👋 downloading...")

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

		/*
			case request.Method == http.MethodPut && authorised == true:
				response.WriteHeader(http.StatusOK)
				response.Write([]byte("🙂 PUT"))

			case request.Method == http.MethodDelete && authorised == true:
				response.WriteHeader(http.StatusOK)
				response.Write([]byte("🙂 DELETE"))
		*/

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}

	}
}
