package server

import (
	"fmt"
	"time"
)

func goRoutineSimplismProcess(currentSimplismProcess SimplismProcess) {
	// TODO: store the process a first time

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// TODO: this is a work in progress 🚧
			fmt.Println("👋", currentSimplismProcess)
			// store data somewhere
			// how to garden the data? (with which condition)
			// create an API endpoint to query the data
		}
	}
}
