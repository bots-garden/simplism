package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	simplismTypes "simplism/types"
	"time"
)

/*
Every Simplism process launched with the --discovery-endpoint flag
will send information to the discovery server.
*/

var delayToSendInformationProcess = 10

// sendProcessInformationToDiscoveryServer sends the current Simplism process information to the discovery server.
//
// Parameters:
// - currentSimplismProcess: The current Simplism process to be sent.
// - wasmArgs: The Wasm arguments containing the discovery endpoint.
func sendProcessInformationToDiscoveryServer(currentSimplismProcess simplismTypes.SimplismProcess, wasmArgs simplismTypes.WasmArguments) {
	discoveryEndPoint := wasmArgs.DiscoveryEndpoint

	// make a post http request
	jsonPayload, err := json.Marshal(currentSimplismProcess)
	if err != nil {
		fmt.Println("😡 Error marshaling JSON:", err)
		return // 🤔
	}

	req, err := http.NewRequest("POST", discoveryEndPoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("😡 Error creating request:", err)
		return // 🤔
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if wasmArgs.AdminDiscoveryToken != "" {
		req.Header.Set("admin-discovery-token", wasmArgs.AdminDiscoveryToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	//fmt.Println("🤖 simplism process:", resp.Body)

	if err != nil {
		fmt.Println("😡 Error sending request:", err)
		return // 🤔
	}
	defer resp.Body.Close()
}

// goRoutineSimplismProcess is a function that processes a given SimplismProcess and WasmArguments in a separate goroutine.
//
// It sends the process information to a discovery server and then repeatedly sends the process information at a predefined interval.
//
// Parameters:
// - currentSimplismProcess: the SimplismProcess to be processed.
// - wasmArgs: the WasmArguments to be used during processing.
func goRoutineSimplismProcess(currentSimplismProcess simplismTypes.SimplismProcess, wasmArgs simplismTypes.WasmArguments) {

	sendProcessInformationToDiscoveryServer(currentSimplismProcess, wasmArgs)

	ticker := time.NewTicker(time.Duration(delayToSendInformationProcess) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//fmt.Println("👋", currentSimplismProcess)
			sendProcessInformationToDiscoveryServer(currentSimplismProcess, wasmArgs)
		}
	}
}
