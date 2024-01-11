package registry

import (
	"net/http"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// getJSONListOfFiles retrieves a JSON list of files from the specified registry path.
//
// Parameters:
// - registryPath: a string representing the path of the registry.
//
// Returns:
// - []byte: a byte slice containing the JSON data of the list of files.
// - error: an error if there was a problem retrieving or converting the list of files.
func getJSONListOfFiles(registryPath string) ([]byte, error) {

	files, err := getListOfFiles(registryPath)
	if err != nil {
		return nil, err
	}

	jsonData, err := convertListOfFilesToJSON(files)
	if err != nil {
		return nil, err
	}

	return jsonData, nil

}

// getTableListOfFiles returns a list of files in the given registry path as a table.
//
// The function takes a registryPath parameter of type string, representing the path to the registry.
// It returns a [][]string, which is a two-dimensional slice containing the name and path of each file in the registry, and an error if any.
func getTableListOfFiles(registryPath string) ([][]string, error) {
	data := [][]string{}

	files, err := getListOfFiles(registryPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data = append(data, []string{
			file.Name,
			file.Path,
			strconv.FormatInt(file.Size, 10),
			file.Created.Format("2006-01-02 15:04:05"),
		})
	}

	return data, nil

}

// getProcessesTableWriter returns a tablewriter.Table initialized with the provided response and data.
//
// The response parameter is an http.ResponseWriter used to write the HTTP response.
// The data parameter is a 2D slice of strings containing the data for the table.
// The function returns a pointer to a tablewriter.Table.
func getListOfFilesTableWriter(response http.ResponseWriter, data [][]string) *tablewriter.Table {
	table := tablewriter.NewWriter(response)
	table.SetHeader([]string{"name", "path", "size", "created"})
	for _, v := range data {
		table.Append(v)
	}

	return table
}
