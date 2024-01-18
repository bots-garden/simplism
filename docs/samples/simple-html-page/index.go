// main package
package main

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/extism/go-pdk"
)

var (
	//go:embed resources/index.html
	html []byte
)

type ResponseData struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Code   int                 `json:"code"`
}

//export handle
func handle() {

	message := "Simplism is propulsed by Extism"

	htmlPage :=
		strings.Replace(
			string(html),
			"{{message}}",
			message,
			-1)

	responseData := ResponseData{
		Body:   htmlPage,
		Header: map[string][]string{"Content-Type": {"text/html; charset=utf-8"}},
		Code:   200,
	}

	jsonResponse, _ := json.Marshal(responseData)

	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)

}

func main() {}
