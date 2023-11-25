// main package
package main

import (
    "encoding/json"
    "github.com/extism/go-pdk"
)

// Argument structure (from the request)
type Argument struct {
    Body   string              `json:"body"`
    Header map[string][]string `json:"header"`
    Method string              `json:"method"`
    URI    string              `json:"uri"`
}

// ReturnValue structure (for the response)
type ReturnValue struct {
    Body   string              `json:"body"`
    Header map[string][]string `json:"header"`
    Code   int                 `json:"code"`
}

//export handle
func handle() {
    // read function argument from the memory
    input := pdk.Input()

    configName, _ := pdk.GetConfig("name")

    pdk.Log(pdk.LogInfo, "⭐️ config name: "+configName)

    /* Expected argument
    type Argument struct {
        Body   string              `json:"body"`
        Header map[string][]string `json:"header"`
        Method string              `json:"method"`
        URI    string              `json:"uri"`
    }
    */
    var argument Argument
    json.Unmarshal(input, &argument)
    pdk.Log(pdk.LogInfo, "📙 content type: "+argument.Header["Content-Type"][0])
    pdk.Log(pdk.LogInfo, "📝 method: "+argument.Method)
    pdk.Log(pdk.LogInfo, "📝 uri:"+argument.URI)
    pdk.Log(pdk.LogInfo, "📝 body:"+argument.Body)
    
    message := "🤗 Hello "
    
    /* Expected response
    type ReturnValue struct {
        Body   string              `json:"body"`
        Header map[string][]string `json:"header"`
        Code   int                 `json:"code"`
    }
    */
    response := ReturnValue{
        Body:   message,
        Header: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
        Code:   200,
    }
    // response to Json string
    jsonResponse, _ := json.Marshal(response)

    // copy output to host memory
    mem := pdk.AllocateBytes(jsonResponse)
    pdk.OutputMemory(mem)

}

func main() {}
