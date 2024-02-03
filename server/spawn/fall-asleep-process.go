package spawn

import (
	"fmt"
	"net/http"
	processesHelper "simplism/helpers/processes"

	"simplism/server/discovery"
	"simplism/server/router"

	//"simplism/server/spawn"
	simplismTypes "simplism/types"
)

// fallAsleepProcess kills a process with the given PID.
//
// It takes an integer parameter `pid` which represents the process ID.
// The function returns a `simplismTypes.SimplismProcess` and an error.
func fallAsleepProcess(pid int) (simplismTypes.SimplismProcess, error) {
	errKill := processesHelper.KillSimplismProcess(pid)
	if errKill != nil {
		fmt.Println("😡 when killing (fall asleep) the process:", errKill)
		return simplismTypes.SimplismProcess{}, errKill
	} else {

		foundProcess, err := discovery.NotifyProcessAsleep(pid)

		// Do not remove the entry from the recovery file
		//delete(spawnedProcesses, foundProcess.HTTPPort)
		//yamlHelper.WriteYamlFile("recovery.yaml", spawnedProcesses)

		// Change the handler
		router.GetRouter().HandleFunc("/service/"+foundProcess.ServiceName, func(response http.ResponseWriter, request *http.Request) {

			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte("[" + foundProcess.HTTPPort + "]🚀(Not found) Simplism process asleep"))
		})

		fmt.Println("😴 Process asleep successfully:", foundProcess.ServiceName)

		if err != nil {
			fmt.Println("😡 handler-spawn/NotifyProcessAsleep", err)
		} else {
			fmt.Println("🙂 Notification for process asleep sent for db update")
		}
		return foundProcess, err
	}

}
