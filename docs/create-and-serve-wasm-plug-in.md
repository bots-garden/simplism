# Create and serve an Extism (wasm) plug-in

Main objectives:
- Create an Extism plug-in
- Serve the plug-in through HTTP like a microservice
- Query the service

> prerequisite: 
> - Install Go and TinyGo compilers
> - Install Simplism

## Request Data

The Simplism server sends the request data to the wasm plug-in in this form:

> example
```json
{
    "body": "hello",
    "method": "POST",
    "uri": "/hello/world",
    "header": {"Content-Type":["text/plain; charset=utf-8"]}
}
```

## Response Data

The wasm plugin has to return a response data in this form:

> example
```json
{
    "body": "hello",
    "header": {"Content-Type":["text/plain; charset=utf-8"]},
    "code": 200
}
```

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

//export handle
func handle() {
    // read request data from the memory
    input := pdk.Input()

    var requestData RequestData
    json.Unmarshal(input, &requestData)
    
    message := "🤗 Hello " + requestData.Body
    
    responseData := ResponseData{
        Body:   message,
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

### Build the wasm plug-in

```bash
tinygo build -scheduler=none --no-debug \
-o hello.wasm \
-target wasi main.go
```

### Serve the wasm plug-in

```bash
simplism listen \
hello.wasm handle --http-port 8080 --log-level info
```
> `handle` is the name of the function to call


### Query the wasm service

```bash
curl http://localhost:8080 \
-d 'Bob Morane'
```
> you should get this response: `🤗 Hello Bob Morane`

### Serve the wasm plug-in with Docker

```bash
docker run \
-p 8080:8080 \
-v $(pwd):/app \
--rm k33g/simplism:0.1.1 \
/simplism listen ./app/hello.wasm handle --http-port 8080 --log-level info
```
