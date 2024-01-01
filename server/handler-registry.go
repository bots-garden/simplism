package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	httpHelper "simplism/helpers/http"
	simplismTypes "simplism/types"
	"strconv"
	"strings"

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

type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	//FileType    string     `json:"fileType"`
}

func fetchFileInfo(directory string) ([]FileInfo, error) {
	var files []FileInfo
	if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() { // Only process files, not directories
			fileInfo := FileInfo{
				Name: info.Name(),
				Path: path,
				//FileType: info.IsDir() == true ? "directory" : "file",
			}
			files = append(files, fileInfo)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func convertFileInfoToJSON(files []FileInfo) (string, error) {
	jsonData, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func registryHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

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

		// download: /registry/pull/{wasmfilename}
		case request.Method == http.MethodGet &&
			strings.HasPrefix(request.RequestURI, "/registry/pull/") &&
			strings.HasSuffix(request.RequestURI, ".wasm") &&
			privateRegistryAuthorised == true:

			//fmt.Println("👋 downloading...", request.RequestURI, request.URL)

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

			files, err := fetchFileInfo(wasmArgs.RegistryPath)
			//fmt.Println(files)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error reading directory"))
				return
			}

			jsonString, err := convertFileInfoToJSON(files)
			//fmt.Println(jsonString)
			if err != nil {
				//fmt.Println("😡 Error writing wasm file")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("😡 Error when converting directory list to JSON"))
				return
			}

			response.WriteHeader(http.StatusOK)
			response.Header().Set("Content-Type", "application/json; charset=utf-8")

			response.Write([]byte(jsonString))

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
