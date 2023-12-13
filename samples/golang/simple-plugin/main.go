package main

import (
	"strconv"
	"github.com/extism/go-pdk"
)

// TextPlainResponse generates a text/plain response body with the provided content and HTTP status code.
//
// Parameters:
//   - bodyContent: the content of the response body as a string.
//   - httpCode: the HTTP status code as an integer.
//
// Returns:
//   - a string representing the text/plain response body.
func TextPlainResponse(bodyContent string, httpCode int) string {
	return `{"code":` + strconv.Itoa(httpCode) + `, "body":"` + bodyContent + `","header":		{"Content-Type":["text/plain; charset=utf-8"]}}`
}

//export say_hello
func say_hello() {
	message := "👋 Hello World 🌍" 
	response := TextPlainResponse(message, 200)
	// copy output to host memory
	mem := pdk.AllocateString(response)
	pdk.OutputMemory(mem)
}

func main() {}
