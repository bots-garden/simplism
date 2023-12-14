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

    /* Expected request data structure
    type RequestData struct {
        Body   string              `json:"body"`
        Header map[string][]string `json:"header"`
        Method string              `json:"method"`
        URI    string              `json:"uri"`
    }
    */
    var requestData RequestData
    json.Unmarshal(input, &requestData)
    pdk.Log(pdk.LogInfo, "📙 content type: "+requestData.Header["Content-Type"][0])
    pdk.Log(pdk.LogInfo, "📝 method: "+requestData.Method)
    pdk.Log(pdk.LogInfo, "📝 uri:"+requestData.URI)
    pdk.Log(pdk.LogInfo, "📝 body:"+requestData.Body)
    
    message := "🤗 Hello "
    
    /* Expected response
    type ResponseData struct {
        Body   string              `json:"body"`
        Header map[string][]string `json:"header"`
        Code   int                 `json:"code"`
    }
    */
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
