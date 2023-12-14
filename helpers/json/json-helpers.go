package jsonHelper

import (
	"encoding/json"
	"fmt"
	simplismTypes "simplism/types"
)

// GetStoreRecordFromJSONString parses a JSON string and returns a simplismTypes.StoreRecord and an error.
//
// It takes a JSON string as a parameter.
// It returns a simplismTypes.StoreRecord and an error.
func GetStoreRecordFromJSONString(json string) (simplismTypes.StoreRecord, error) {
    return GetStoreRecordFromJSONBytes([]byte(json))
}

// GetStoreRecordFromJSONBytes parses the given JSON bytes and returns a StoreRecord and an error.
//
// The function takes a byte array 'body' as input, which contains the JSON data to be parsed.
// It returns a StoreRecord and an error. The StoreRecord is the parsed JSON data, and the error
// indicates if there was any issue during the parsing process.
func GetStoreRecordFromJSONBytes(body []byte) (simplismTypes.StoreRecord, error) {
	var record simplismTypes.StoreRecord
	jsonUnmarshallErr := json.Unmarshal(body, &record)
	if jsonUnmarshallErr != nil {
		fmt.Println("😡 Error when unmarshaling JSON:", jsonUnmarshallErr)
	}
    return record, jsonUnmarshallErr
}

// GetSimplismProcesseFromJSONString generates a SimplismProcess object from a JSON string.
//
// json: the JSON string from which to generate the SimplismProcess object.
// Returns: a SimplismProcess object and an error if there was a problem generating it.
func GetSimplismProcesseFromJSONString(json string) (simplismTypes.SimplismProcess, error) {
    return GetSimplismProcesseFromJSONBytes([]byte(json))
}

// GetSimplismProcesseFromJSONBytes retrieves a SimplismProcess object from a JSON byte array.
//
// The function takes a byte array `body` as input.
// It returns a SimplismProcess object and an error.
func GetSimplismProcesseFromJSONBytes(body []byte) (simplismTypes.SimplismProcess, error) {
	var simplismProcess simplismTypes.SimplismProcess
	jsonUnmarshallErr := json.Unmarshal(body, &simplismProcess)
	if jsonUnmarshallErr != nil {
		fmt.Println("😡 Error when unmarshaling JSON:", jsonUnmarshallErr)
	}
    return simplismProcess, jsonUnmarshallErr
}

// GetJSONBytesFromSimplismProcess returns the JSON representation of the given simplismProcess and any error encountered.
//
// It takes a simplismProcess of type simplismTypes.SimplismProcess as the parameter.
// It returns a byte slice ([]byte) containing the JSON representation of the simplismProcess and an error.
func GetJSONBytesFromSimplismProcess(simplismProcess simplismTypes.SimplismProcess) ([]byte, error) {
    jsonProcess, jsonMarshallErr := json.Marshal(simplismProcess)
    if jsonMarshallErr != nil {
        fmt.Println("😡 Errorwhen  marshaling JSON:", jsonMarshallErr)
    }
    return jsonProcess, jsonMarshallErr
}
