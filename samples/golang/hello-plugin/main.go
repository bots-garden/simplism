// main package
package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
)

// Human structure
type Human struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

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

//export say_hello
func say_hello() {
	// read function argument from the memory
	input := pdk.Input()

	/* Expected
	type Argument struct {
		Body   string              `json:"body"`
		Header map[string][]string `json:"header"`
		Method string              `json:"method"`
		URI    string              `json:"uri"`
	}
	*/
	var argument Argument
	json.Unmarshal(input, &argument)

	// ✋ displaying messages slows down the plugin execution
	pdk.Log(pdk.LogInfo, "📙 content type: "+argument.Header["Content-Type"][0])
	pdk.Log(pdk.LogInfo, "📝 method: "+argument.Method)
	pdk.Log(pdk.LogInfo, "📝 uri:"+argument.URI)
	
	var message string
	var code = 200
	var human Human
	errHuman := json.Unmarshal([]byte(argument.Body), &human)
	if errHuman != nil {
		message = "😡 Hello John Doe"
		code = 500
	} else {
		message = "🤗 Hello " + human.FirstName + " " + human.LastName
	}

	
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
		Code:   code,
	}
	// response to Json string
	jsonResponse, _ := json.Marshal(response)

	// copy output to host memory
	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)

}

func main() {}
