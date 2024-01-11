package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// setHeaders represents a function that sets the headers for a given HTTP response.
//
// Parameters:
// - response: the HTTP response writer.
// - name: the filename for the attachment.
// - len: the length of the data.
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

// getListOfFiles returns a list of file information for all the files in the specified directory.
//
// directory: the directory path to search for files
// []FileInfo: a slice of FileInfo struct representing the file information
// error: an error if any occurred during the file search process
func getListOfFiles(directory string) ([]FileInfo, error) {
	var files []FileInfo
	if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
        
		if !info.IsDir() { // Only process files, not directories
			fileInfo := FileInfo{
				Name: info.Name(),
				Path: path,
                Created: info.ModTime(),
                Size: info.Size(),
				//FileType: info.IsDir() == true ? "directory" : "file",
			}
            // add only wasm files
            if strings.HasSuffix(fileInfo.Path, ".wasm") {
                files = append(files, fileInfo)
            }
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

// convertListOfFilesToJSON converts a list of FileInfo objects to JSON format.
//
// files - The list of FileInfo objects to be converted.
// []byte - The converted JSON data.
// error - An error, if any occurred during the conversion.
func convertListOfFilesToJSON(files []FileInfo) ([]byte, error) {
	jsonData, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
