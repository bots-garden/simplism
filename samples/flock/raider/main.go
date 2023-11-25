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

func displayConfigInfo(raiderName, basestarURL, httpPort string) {
	pdk.Log(pdk.LogInfo, "🚀 raider name: "+raiderName)   // name of the raider
	pdk.Log(pdk.LogInfo, "🌍 basestar URL: "+basestarURL) // url of the basestar
	pdk.Log(pdk.LogInfo, "🌎 http port: "+httpPort)
}

func displayReceivedData(input []byte) {
	// Received data
	var argument Argument
	json.Unmarshal(input, &argument)
	pdk.Log(pdk.LogInfo, "📙 content type: "+argument.Header["Content-Type"][0])
	pdk.Log(pdk.LogInfo, "📝 method: "+argument.Method)
	pdk.Log(pdk.LogInfo, "📝 uri:"+argument.URI)
	pdk.Log(pdk.LogInfo, "📝 body:"+argument.Body)
}

func sendResponse(message string) {
	response := ReturnValue{
		Body:   `{"message": "` + message + `"}`,
		Header: map[string][]string{"Content-Type": {"application/json; charset=utf-8"}},
		Code:   200,
	}
	// response to Json string
	jsonResponse, _ := json.Marshal(response)

	// copy output to host memory
	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)

}

//export handle
func handle() {
	// read function argument from the memory
    /*
	input := pdk.Input()
    displayReceivedData(input)
    */

	// get config information
	raiderName, _ := pdk.GetConfig("name")
	basestarURL, _ := pdk.GetConfig("basestar")
	httpPort, _ := pdk.GetConfig("http-port")

	//displayConfigInfo(raiderName, basestarURL, httpPort)

	// send request to the basestar
	req := pdk.NewHTTPRequest("POST", basestarURL)
	req.SetHeader("Content-Type", "application/json")
	req.SetBody([]byte(`{"name":"` + raiderName + `", "url":"http://localhost:` + httpPort + `"}`))
	res := req.Send()

	// display response from the basestar
	pdk.Log(pdk.LogInfo, "📩 from basestar:"+string(res.Body()))

	sendResponse("🚀 Hello I am " + raiderName)

}

func main() {}
