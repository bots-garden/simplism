package discovery

import (
	"encoding/json"
	"net/http"
	data "simplism/server/data"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"go.etcd.io/bbolt"
)

// getJSONProcesses retrieves the list of services that are running from the database and returns it as a JSON byte array.
//
// Parameters:
// - db: a pointer to a bbolt.DB instance representing the database.
//
// Returns:
// - []byte: the JSON byte array containing the list of services.
// - error: an error if there was a problem retrieving the processes or encoding them as JSON.
func getJSONProcesses(db *bbolt.DB) ([]byte, error) {
	// get the list of the services that are running
	processes, dbErr := data.GetSimplismProcessesListFromDB(db)
	if dbErr != nil {
		return nil, dbErr
	} 
	jsonBuffer, jsoneErr := json.Marshal(processes)

	if jsoneErr != nil {
		return nil, jsoneErr
	} else {
		return jsonBuffer, nil
	}
}

// getTableProcesses is a function that retrieves the list of processes from the provided bbolt.DB.
//
// It takes a single parameter:
// - db: a pointer to a bbolt.DB instance.
//
// It returns a 2D slice of strings and an error:
// - [][]string: the list of processes, where each process is represented by a slice of strings.
// - error: any error that occurred during the retrieval process.
func getTableProcesses(db *bbolt.DB) ([][]string, error) {
	processes, dbErr := data.GetSimplismProcessesListFromDB(db)
	data := [][]string{}

    if dbErr != nil {
        return nil, dbErr
    }

	for _, process := range processes {

		
		isStopped := func() string {
			if process.StopTime.IsZero() {
				return ""
			} else {
				return "x"
			}
		}

		isStarted := func() string {
			if process.StartTime.IsZero() {
				return ""
			} else {
				if process.StopTime.IsZero() != true {
					return ""
				}
				return "x"
			}
		}
		

		data = append(data, []string{
			strconv.Itoa(process.PID),
			process.ServiceName,
			process.HTTPPort,
			process.FunctionName,
			process.FilePath,
			//process.Information,
			isStarted(),
			isStopped(),
			//process.StartTime.Format("2006-01-02 15:04:05"),
			//process.StopTime.Format("2006-01-02 15:04:05"),

			//process.StopTime.Local().String(),
			//process.RecordTime.Local().String(),

		})
	}
	return data, nil
}

// getProcessesTableWriter returns a tablewriter.Table that displays the provided data in a tabular format.
//
// The function takes two parameters:
// - response: an http.ResponseWriter object used to write the HTTP response.
// - data: a 2D slice of strings representing the data to be displayed in the table.
//
// The function returns a pointer to a tablewriter.Table object.
func getProcessesTableWriter(response http.ResponseWriter, data [][]string) *tablewriter.Table {
	table := tablewriter.NewWriter(response)
	table.SetHeader([]string{"pid", "name", "http", "function", "path", "started", "stopped"})
	for _, v := range data {
		table.Append(v)
	}

	return table
}
