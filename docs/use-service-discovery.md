# Use the service discovery feature

You can use a Simplism service as a service discovery service.

## Create a simple Simplism service

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

//export handle
func handle() {    
    message := "👋 I'm the discovery service"
    
    responseData := ResponseData{
        Body:   message,
        Header: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
        Code:   200,
    }
    jsonResponse, _ := json.Marshal(responseData)

    mem := pdk.AllocateBytes(jsonResponse)
    pdk.OutputMemory(mem)

}

func main() {}
```

## Build it

```bash
tinygo build -scheduler=none --no-debug \
-o service-discovery.wasm \
-target wasi main.go
```

## Run the service as a discovery service

To activate the discovery feature, use these flags:

- `--service-discovery true`
- `--admin-discovery-token <token>`

> start the discovery service with the following command:
```bash
simplism listen \
service-discovery.wasm handle \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange
```

## Make the other services discoverable

You can make other services discoverable by using these flags:

- `--discovery-endpoint <endpoint>`
- `--service-name <name>`
- `--admin-discovery-token <token>`

> start the other services with the following commands:
```bash
simplism listen \
../hello-people/hello-people.wasm handle \
--http-port 8081 \
--log-level info \
--service-name hello-people \
--admin-discovery-token people-are-strange \
--discovery-endpoint http://localhost:9000/discovery

simplism listen \
../hello/hello.wasm handle \
--http-port 8082 \
--log-level info \
--service-name hello \
--admin-discovery-token people-are-strange \
--discovery-endpoint http://localhost:9000/discovery
```

## Get the service list

```bash
curl http://localhost:9000/discovery \
-H 'admin-discovery-token:people-are-strange'
```

You should get:
```json
{"10504":{"pid":10504,"functionName":"handle","filePath":"../hello-people/hello-people.wasm","recordTime":"2023-12-13T13:23:08.152474368Z","startTime":"2023-12-13T13:22:58.042558546Z","stopTime":"0001-01-01T00:00:00Z","httpPort":"8081","information":"","serviceName":"hello-people","host":""},"10551":{"pid":10551,"functionName":"handle","filePath":"../hello/hello.wasm","recordTime":"2023-12-13T13:23:03.587034276Z","startTime":"2023-12-13T13:23:03.495780517Z","stopTime":"0001-01-01T00:00:00Z","httpPort":"8082","information":"","serviceName":"hello","host":""}
```

## Query the services

```bash
curl http://localhost:8081 \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'

curl http://localhost:8082 \
-d 'Bob Morane'
```

## Query the services through the discovery service

As soon as service discovery is enabled, you can use it directly to query services by name:

```bash
curl http://localhost:9000/service/hello-people \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'

curl http://localhost:9000/service/hello \
-d 'Bob Morane'
```

