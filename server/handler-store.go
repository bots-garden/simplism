package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	httpHelper "simplism/helpers/http"
	jsonHelper "simplism/helpers/json"
	simplismTypes "simplism/types"
)

func storeHandler(wasmArgs simplismTypes.WasmArguments) http.HandlerFunc {

	db, err := initializeStoreDB(wasmArgs, "") // use the default path
	if err != nil {
		panic(err) // TODO: handle error in a better way
	} else {
		fmt.Println("📦 Store database initialized")
	}

	return func(response http.ResponseWriter, request *http.Request) {

		authorised := httpHelper.CheckStoreToken(request, wasmArgs)

		switch { // /store
		case request.Method == http.MethodPost && authorised == true:
			body := httpHelper.GetBody(request)
			record, err := jsonHelper.GetStoreRecordFromJSONBytes(body)

			if err != nil {
				fmt.Println("😡 Error when unmarshalling record:", err)
				response.WriteHeader(http.StatusInternalServerError)
			} else {

				err = saveToStore(db, record.Key, record.Value)
				if err != nil {
					fmt.Println("😡 Error when saving record to DB:", err)
					response.WriteHeader(http.StatusInternalServerError)

				} else {
					response.WriteHeader(http.StatusOK)
					response.Write([]byte("🙂 record with key " + record.Key + " saved"))
				}
			}
		
		case request.Method == http.MethodGet && authorised == true:
			// all; http://localhost:8080/store
			// prefix: http://localhost:8080/store?prefix=hello
			query := request.URL.Query()

			keyList, present := query["key"]
			if present {

				key := keyList[0]
				record := getFromStore(db, key)
				jsonString, err := json.Marshal(record)
				if err != nil {
					fmt.Println("😡 When marshalling", err)
					response.WriteHeader(http.StatusInternalServerError)
				} else {
					response.WriteHeader(http.StatusOK)
					response.Write(jsonString)
				}

			} else {

				prefixList, present := query["prefix"]
				if !present {
					// get all records
					records := getAllFromStore(db)
					jsonString, err := json.Marshal(records)
					if err != nil {
						fmt.Println("😡 When marshalling", err)
						response.WriteHeader(http.StatusInternalServerError)
					} else {
						response.WriteHeader(http.StatusOK)
						response.Write(jsonString)
					}
	
				} else {
					// get records with prefix
					prefix := prefixList[0]
					records := getAllWitPrefixFromStore(db, prefix)
					jsonString, err := json.Marshal(records)
	
					if err != nil {
						fmt.Println("😡 When marshalling", err)
						response.WriteHeader(http.StatusInternalServerError)
					} else {
						response.WriteHeader(http.StatusOK)
						response.Write(jsonString)
					}
	
				}
			}



			//response.WriteHeader(http.StatusOK)
			//response.Write([]byte("📦 Hello [GET]"))

		case request.Method == http.MethodPut && authorised == true:
			response.WriteHeader(http.StatusOK)
			response.Write([]byte("📦 Hello [PUT]"))

		case request.Method == http.MethodDelete && authorised == true:
			// key: http://localhost:8080/store?key=abcd
			query := request.URL.Query()

			keyList, present := query["key"]
			if !present {
				response.WriteHeader(http.StatusNotFound)
				response.Write([]byte("key not present"))
			} else {

				key := keyList[0]
				err := deleteFromStore(db, key)
				if err != nil {
					fmt.Println("😡 When deleting", err)
					response.WriteHeader(http.StatusInternalServerError)
				} else {
					response.WriteHeader(http.StatusOK)
					response.Write([]byte("🙂 record with key " + key + " deleted"))
				}
			}

		case authorised == false:
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("😡 You're not authorized"))

		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte("😡 Method not allowed"))
		}
	}
}
