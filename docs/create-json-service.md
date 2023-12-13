# Create a JSON service with an Extism plug-in


## Create an Extism plugin

### Create the source code
```bash
go mod init hello
touch main.go
```

> `main.go`
```go
// main package
package main

import (
    "encoding/json"
    "github.com/extism/go-pdk"
)

// RequestData structure (from the request)
type RequestData struct {
    Body   string              `json:"body"`
    Header map[string][]string `json:"header"`
    Method string              `json:"method"`
    URI    string              `json:"uri"`
}

// ResponseData structure (for the response)
type ResponseData struct {
    Body   string              `json:"body"`
    Header map[string][]string `json:"header"`
    Code   int                 `json:"code"`
}

type Human struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
}

//export handle
func handle() {
    // read request data from the memory
    input := pdk.Input()

    var requestData RequestData
    json.Unmarshal(input, &requestData)

    var human Human
    json.Unmarshal([]byte(requestData.Body), &human)
    
    message := "🤗 Hello " + human.FirstName + " " + human.LastName
    
    responseData := ResponseData{
        Body:   `{"message": "` + message + `"}`,
        Header: map[string][]string{"Content-Type": {"application/json; charset=utf-8"}},
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

### Build the wasm plug-in

```bash
tinygo build -scheduler=none --no-debug \
-o hello-people.wasm \
-target wasi main.go
```

### Serve the wasm plug-in

```bash
simplism listen \
hello-people.wasm handle --http-port 8080 --log-level info
```
> `handle` is the name of the function to call


### Query the wasm service

```bash
curl http://localhost:8080 -v \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```
> you should get this response: `{"message": "🤗 Hello Bob Morane"}`
