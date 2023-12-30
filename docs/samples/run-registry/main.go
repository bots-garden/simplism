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

	responseData := ResponseData{
		Body:   "🐳 tiny wasm registry 🦄, use the `/registry` endpoint",
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
