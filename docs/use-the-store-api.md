# Use the store API

A Simplism process/server can expose a "store HTTP API" that can be used to save/read data with this form : `{"key":"","value":""}`.

The other Simplism process/server can use this API to read, save and delete data from the store.

You can of course start several stores.

> The data are stored thanks to the **[BBolt project](https://github.com/etcd-io/bbolt)**, an embedded key/value database for Go.

## Start a Simplism store

You need to use the `--store-mode true` flag to enable the store mode. And you can protect the store access with the `--admin-store-token` flag.

```bash
simplism listen \
store.wasm handle \
--http-port 8080 \
--log-level info \
--store-mode true \
--admin-store-token morrison-hotel \
--information "👋 I'm the store service" \
--service-name store
```
> You can define the path of the store with the `--store-path` flag. Otherwise the store db file will be stored in the same folder as the wasm file (and the name will be the wasm file name + ".store.db")
This will activate a new endpoint at `http://localhost:8080/store`.

It's possible to use a configuration file to start the Simplism store:

```yaml
store-config:
  wasm-file: ./store.wasm
  wasm-function: handle
  http-port: 8080
  log-level: info
  store-mode: true
  admin-store-token: morrison-hotel
  information: 👋 I'm the store service
  service-name: store
```
> You can define the path of the store with the `store-path` field. Otherwise the store db file will be stored in the same folder as the wasm file (and the name will be the wasm file name + ".store.db")

## Make queries to the store

### Add data to the store

```bash
curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"hello","value":"hello world"}'
```

### Read a record from the store

```bash
curl http://localhost:8080/store?key=hello \
-H 'admin-store-token: morrison-hotel'
```

### Read all the records from the store

```bash
curl http://localhost:8080/store \
-H 'admin-store-token: morrison-hotel'
```

### Get all records with a key starting with a prefix

```bash
curl http://localhost:8080/store?prefix=hel \
-H 'admin-store-token: morrison-hotel'
```

### Delete a record from the store

```bash
curl -X DELETE http://localhost:8080/store?key=hello \
-H 'admin-store-token: morrison-hotel'
```

## Make queries to the store from another Simplism process

Create a new Simplism service with this source code:

```go
// main package
package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

// ResponseData structure (for the response)
type ResponseData struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Code   int                 `json:"code"`
}

func addRecord(storeURL string, storeToken string, key string, value string) {
	pdk.Log(pdk.LogInfo, "🌍 store-url: "+storeURL)

	req := pdk.NewHTTPRequest("POST", storeURL+"/store")
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("admin-store-token", storeToken)
	req.SetBody([]byte(`{"key":"` + key + `","value":"` + value + `"}`))
	res := req.Send()

	pdk.Log(pdk.LogInfo, "📦 store response: "+string(res.Body()))
}

//export handle
func handle() {

	storeURL, _ := pdk.GetConfig("store-url")
	storeToken, _ := pdk.GetConfig("store-token")

	addRecord(storeURL, storeToken, "001", "first record")
	addRecord(storeURL, storeToken, "002", "second record")
	addRecord(storeURL, storeToken, "003", "third record")

	responseData := ResponseData{
		Body:   "records added",
		Header: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		Code:   200,
	}
	// response to Json string
	jsonResponse, _ := json.Marshal(responseData)

	// copy output to host memory
	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)

}

func main() {}
```

### Start the service

Start the service with this configuration file:

```yaml
use-store-config:
  wasm-file: ./use-store.wasm
  wasm-function: handle
  http-port: 9090
  log-level: info
  config: |
    {
      "store-url":"http://localhost:8080",
      "store-token":"morrison-hotel"
    }
  service-name: use-store
```

Then, do this HTTP request:

```bash
curl http://localhost:9090
```
This will add three records to the store. And you can now query them from the store with this request:

```bash
curl http://localhost:8080/store \
-H 'admin-store-token: morrison-hotel'
```

You should get the following output:
```json
{"001":"first record","002":"second record","003":"third record"}
```
