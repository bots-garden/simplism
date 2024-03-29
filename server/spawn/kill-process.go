package spawn

import (
	"fmt"
	"net/http"
	processesHelper "simplism/helpers/processes"
	yamlHelper "simplism/helpers/yaml"

	"simplism/server/discovery"
	"simplism/server/router"
	simplismTypes "simplism/types"
)

// killProcess kills a process with the given PID.
//
// It takes an integer parameter `pid` which represents the process ID.
// The function returns a `simplismTypes.SimplismProcess` and an error.
func killProcess(pid int) (simplismTypes.SimplismProcess, error) {
	errKill := processesHelper.KillSimplismProcess(pid)
	if errKill != nil {
		fmt.Println("😡 when killing the process:", errKill)
		return simplismTypes.SimplismProcess{}, errKill
	} else {

		foundProcess, err := discovery.NotifyProcessKilled(pid)

		// Update the recovery file (remove the entry for the killed process) from the map
		delete(spawnedProcesses, foundProcess.HTTPPort)

		yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)

		// Change the handler

		router.GetRouter().HandleFunc("/service/"+foundProcess.ServiceName, func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte("(Not found) Simplism process killed"))
		})

		fmt.Println("🙂 Process killed successfully:", foundProcess.ServiceName)

		if err != nil {
			fmt.Println("😡 handler-spawn/NotifyProcessKilled", err)
		} else {
			fmt.Println("🙂 Notification for process killed sent for db update")
		}
		return foundProcess, err
	}

}
