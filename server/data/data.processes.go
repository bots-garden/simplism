package data

import (
	"fmt"
	"log"
	jsonhelper "simplism/helpers/json"
	simplismTypes "simplism/types"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

// initializeProcessesDB initializes the database for the given WasmArguments.
//
// It takes a single parameter, wasmArgs, of type simplismTypes.WasmArguments.
// It returns a *bolt.DB and an error.
func InitializeProcessesDB(wasmArgs simplismTypes.WasmArguments) (*bolt.DB, error) {

	/*
		go install go.etcd.io/bbolt/cmd/bbolt@latest
		bbolt keys samples/flock/discovery/discovery.wasm.db simplism-bucket
		bbolt page --all samples/flock/discovery/discovery.wasm.db
	*/

	//db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	db, err := bolt.Open(wasmArgs.FilePath+".processes.db", 0600, nil)
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
func SaveSimplismProcessToDB(db *bolt.DB, simplismProcess simplismTypes.SimplismProcess) error {
	simplismProcess.RecordTime = time.Now()
	// convert PID to string
	pidStr := strconv.Itoa(simplismProcess.PID)

	// TODO: make some investigations: what if key == "key-service_name"

	// convert the process information to JSON
	jsonProcess, _ := jsonhelper.GetJSONBytesFromSimplismProcess(simplismProcess)

	// for debugging (temporary)
	//fmt.Println("🟣", string(jsonProcess))

	// Store the process information
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-bucket"))
		err := b.Put([]byte(pidStr), jsonProcess)
		return err
	})
	return err
}

func GetSimplismProcessByPiD(db *bolt.DB, pid int) simplismTypes.SimplismProcess {

	var simplismProcess simplismTypes.SimplismProcess

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-bucket"))
		processValue := b.Get([]byte(strconv.Itoa(pid)))

		if processValue == nil {
			simplismProcess = simplismTypes.SimplismProcess{}
		}
		simplismProcess, _ = jsonhelper.GetSimplismProcesseFromJSONBytes(processValue)
		return nil
	})
	// TODO handle error return value
	return simplismProcess // if nil, return an empty simplismProcess
}

func DeleteSimplismProcessByPiD(db *bolt.DB, pid int) error {

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-bucket"))
		err := b.Delete([]byte(strconv.Itoa(pid)))
		return err
	})
	return err
}

func GetSimplismProcessByName(db *bolt.DB, serviceName string) (simplismTypes.SimplismProcess) {

	var simplismProcess simplismTypes.SimplismProcess

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-bucket"))

		c := b.Cursor()
		for pid, processValue := c.First(); pid != nil; pid, processValue = c.Next() {
			simplismProcess, _ = jsonhelper.GetSimplismProcesseFromJSONBytes(processValue)

			if simplismProcess.ServiceName == serviceName {
				return nil
			}
		}
		return nil
	})
	// TODO handle error return value
	return simplismProcess // if nil, return an empty simplismProcess
}


func GetSimplismProcessesListFromDB(db *bolt.DB) (map[string]simplismTypes.SimplismProcess, error) { // map[string]simplismTypes.SimplismProcess {
	processes := map[string]simplismTypes.SimplismProcess{}

	dbErr := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("simplism-bucket"))

		c := b.Cursor()

		for pid, processValue := c.First(); pid != nil; pid, processValue = c.Next() {
			//fmt.Printf("key=%s, value=%s\n", pid, processValue)
			simplismProcess, _ := jsonhelper.GetSimplismProcesseFromJSONBytes(processValue)
			processes[string(pid)] = simplismProcess
		}

		return nil
	})
	return processes, dbErr
}
