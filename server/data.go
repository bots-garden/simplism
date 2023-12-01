package server

import (
	"fmt"
	"log"
	"simplism/jsonhelper"
	simplismTypes "simplism/types"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

var currentSimplismProcess = simplismTypes.SimplismProcess{}

// initializeDB initializes the database for the given WasmArguments.
//
// It takes a single parameter, wasmArgs, of type simplismTypes.WasmArguments.
// It returns a *bolt.DB and an error.
func initializeDB(wasmArgs simplismTypes.WasmArguments) (*bolt.DB, error) {

	/*
		go install go.etcd.io/bbolt/cmd/bbolt@latest
		bbolt keys samples/flock/discovery/discovery.wasm.db simplism-bucket
		bbolt page --all samples/flock/discovery/discovery.wasm.db
	*/

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
	return db, err
}

// saveSimplismProcessToDB saves the simplism process to the database.
//
// It takes the following parameter(s):
// - db: a pointer to the bolt DB instance.
// - simplismProcess: the simplism process to be saved.
//
// It returns an error.
func saveSimplismProcessToDB(db *bolt.DB, simplismProcess simplismTypes.SimplismProcess) error {
	simplismProcess.RecordTime = time.Now()
	// convert PID to string
	pidStr := strconv.Itoa(simplismProcess.PID)
	// convert the process information to JSON
	jsonProcess, _ := jsonhelper.GetJSONBytesFromSimplismProcess(simplismProcess)
	// for debugging (temporary)
	fmt.Println("🟣", string(jsonProcess))

	// Store the process information
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-bucket"))
		err := b.Put([]byte(pidStr), jsonProcess)
		return err
	})
    return err
}
