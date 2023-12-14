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

func addRecord(storeURL string, storeToken string, key string, value string) {
	pdk.Log(pdk.LogInfo, "🌍 store-url: "+storeURL)

	req := pdk.NewHTTPRequest("POST", storeURL+"/store")
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("admin-store-token", storeToken)
	req.SetBody([]byte(`{"key":"` + key + `","value":"` + value + `"}`))
	res := req.Send()

	pdk.Log(pdk.LogInfo, "📦 store response: "+string(res.Body()))
}

//export handle
func handle() {

	storeURL, _ := pdk.GetConfig("store-url")
	storeToken, _ := pdk.GetConfig("store-token")

	addRecord(storeURL, storeToken, "001", "first record")
	addRecord(storeURL, storeToken, "002", "second record")
	addRecord(storeURL, storeToken, "003", "third record")

	responseData := ResponseData{
		Body:   "records added",
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
