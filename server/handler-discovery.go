package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simplism/httphelper"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

func checkDiscoveryToken(request *http.Request, wasmArgs WasmArguments) bool {
	var authorised bool = false
	// read the header admin-discovery-token
	adminDiscoveryToken := request.Header.Get("admin-discovery-token")

	if wasmArgs.AdminDiscoveryToken != "" {
		// token is awaited
		if wasmArgs.AdminDiscoveryToken == adminDiscoveryToken {
			authorised = true
		} else {
			authorised = false
		}
	} else {
		// check if the env variable ADMIN_DISCOVERY_TOKEN is set
		enAdminDiscoveryToken := os.Getenv("ADMIN_DISCOVERY_TOKEN")
		if enAdminDiscoveryToken != "" {
			// token is awaited
			if enAdminDiscoveryToken == adminDiscoveryToken {
				authorised = true
			} else {
				authorised = false
			}
		} else {
			authorised = true
		}

	}

	return authorised
}

func discoveryHandler(wasmArgs WasmArguments) http.HandlerFunc {
	fmt.Println("🔎 discovery mode activated: /discovery  (", wasmArgs.HTTPPort, ")")

	//db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	db, err := bolt.Open(wasmArgs.FilePath+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("simplism-bucket"))
		if err != nil {
			return fmt.Errorf("😡 When creating bucket: %s", err)
		}
		return nil
	})

	/*
		go install go.etcd.io/bbolt/cmd/bbolt@latest
		bbolt keys samples/flock/discovery/discovery.wasm.db simplism-bucket
		bbolt page --all samples/flock/discovery/discovery.wasm.db
	*/

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := checkDiscoveryToken(request, wasmArgs)
		// Test if it is a POST request
		if request.Method == http.MethodPost && authorised == true {

			body := httphelper.GetBody(request) // process information from simplism POST request

			// create SimpleProcess struct instance from JSON Body
			var simplismProcess SimplismProcess
			jsonUnmarshallErr := json.Unmarshal(body, &simplismProcess)
			if jsonUnmarshallErr != nil {
				fmt.Println("😡 Error when unmarshaling JSON:", jsonUnmarshallErr)
			}
			// record the time
			simplismProcess.RecordTime = time.Now()
			// convert PID to string
			pidStr := strconv.Itoa(simplismProcess.PID)

			// convert the process information to JSON
			jsonProcess, jsonMarshallErr := json.Marshal(simplismProcess)
			if jsonMarshallErr != nil {
				fmt.Println("😡 Errorwhen  marshaling JSON:", err)
			}

			// for debugging (temporary)
			fmt.Println("🟣", string(jsonProcess))

			// Store the process information
			err := db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("simplism-bucket"))
				err := b.Put([]byte(pidStr), jsonProcess)
				return err
			})
			// TODO: look at old records and delete old ones
			// TODO: move all the db stuff to data.go
			if err != nil {
				fmt.Println("😡 When updating bucket", err)
				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.WriteHeader(http.StatusOK)
			}

		} else {
			if authorised == false {
				response.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(response, "😡 You're not authorized")

			} else {
				response.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintln(response, "😡 Method not allowed")
			}
			// 🚧 This is a Work In Progress
			// If GET request

			// If DELETE request

			// If PUT request
		}

	}

}
