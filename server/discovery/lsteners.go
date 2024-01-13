package discovery

import (
	"errors"
	"fmt"
	data "simplism/server/data"
	simplismTypes "simplism/types"

	"go.etcd.io/bbolt"
)

func processKilledListener(db *bbolt.DB) func(pid int) (simplismTypes.SimplismProcess, error) {
	// This function is called by the spawn handler (DELETE method), see handle-spawn.go
	notifyProcessKilled := func(pid int) (simplismTypes.SimplismProcess, error) {
		simplismProcess := data.GetSimplismProcessByPiD(db, pid)

		if simplismProcess.PID == 0 {
			return simplismTypes.SimplismProcess{}, errors.New("😡 Process not found")
		} else {

			// delete from the memory map
			delete(wasmFunctionHandlerList, simplismProcess.ServiceName)
			// delete from database
			err := data.DeleteSimplismProcessByPiD(db, pid)
			if err != nil {
				fmt.Println("😡 When updating bucket for process deletions", err)
				return simplismTypes.SimplismProcess{}, err

			} else {

				fmt.Println("🙂 Bucket updated: process deleted")
			}
			return simplismProcess, nil
		}

	}
	return notifyProcessKilled
}

func processInformationListener(db *bbolt.DB) func(serviceName string) (simplismTypes.SimplismProcess, error) {
	// This function is called by the spawn handler (DELETE method), see handle-spawn.go
	// When we wan to kill a process by it's name
	notifyProcessInformation := func(serviceName string) (simplismTypes.SimplismProcess, error) {

		simplismProcess := data.GetSimplismProcessByName(db, serviceName)

		if simplismProcess.PID == 0 {
			return simplismTypes.SimplismProcess{}, errors.New("😡 Process not found")
		}
		return simplismProcess, nil
	}
    return notifyProcessInformation
}
