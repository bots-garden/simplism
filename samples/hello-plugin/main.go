package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
)

type Human struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type MainArgument struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Method string              `json:"method"`
	Uri    string              `json:"uri"`
}

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
	type MainArgument struct {
		Body   string              `json:"body"`
		Header map[string][]string `json:"header"`
		Method string              `json:"method"`
		Uri    string              `json:"uri"`
	}
	*/
	var argument MainArgument
	errArg := json.Unmarshal(input, &argument)
	if errArg != nil {
		// handle error
	}

	var message = ""
	if argument.Method == "POST" {

		var human Human
		errHuman := json.Unmarshal([]byte(argument.Body), &human)
		if errHuman != nil {
			// handle error😉
		}

		message = "🤗 Hello " + human.FirstName + " " + human.LastName
	} else {
		message = "👋 Hello " 
	}

	/*
		pdk.Log(pdk.LogInfo, "📝 method: "+method)
		pdk.Log(pdk.LogInfo, "📝 uri:"+uri)
		pdk.Log(pdk.LogInfo, "📝 content type:"+contentType)
	*/

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
	jsonResponse, errResponse := json.Marshal(response)
	if errResponse != nil {
		// handle error
	}

	// copy output to host memory
	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)

}

func main() {}
