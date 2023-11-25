// main package
package main

import (
	"encoding/json"
	"strconv"

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

type RaiderMessage struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func displayRequestContent(argument RequestData) {
	pdk.Log(pdk.LogInfo, "⭐️📙 content type: "+argument.Header["Content-Type"][0])
	pdk.Log(pdk.LogInfo, "⭐️📝 method: "+argument.Method)
	pdk.Log(pdk.LogInfo, "⭐️📝 uri:"+argument.URI)
	pdk.Log(pdk.LogInfo, "⭐️📝 body:"+argument.Body)
}

func displayBasestarInfo() {
	basestarName, _ := pdk.GetConfig("name")
	pdk.Log(pdk.LogInfo, "⭐️🛠️ basestar name: "+basestarName)
}

func displayCounter() {
	count := pdk.GetVarInt("count")
	pdk.Log(pdk.LogInfo, "---------------------------------------------")
	pdk.Log(pdk.LogInfo, "⭐️🤖 counter: "+strconv.Itoa(count))
	pdk.Log(pdk.LogInfo, "---------------------------------------------")
}

func getRequestData() RequestData {
	// read function argument from the memory
	requestDataBytes := pdk.Input()

	var requestData RequestData // this is the data request
	json.Unmarshal(requestDataBytes, &requestData)
	return requestData
}

func getRaiderMessage(requestData RequestData) RaiderMessage {
	// unmarshall and get the value of the raider message
	var raiderMessage RaiderMessage
	json.Unmarshal([]byte(requestData.Body), &raiderMessage)
	return raiderMessage
}

func sendResponse(message string) {
	// send response to the raider
	response := ResponseData{
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
	var count int
	count = pdk.GetVarInt("count")
	count = count + 1
	pdk.SetVarInt("count", count)
	
    displayCounter()

	// displayBasestarInfo()

	requestData := getRequestData()
	raiderMessage := getRaiderMessage(requestData)

    sendResponse("🤗 Hello ⭐️" + raiderMessage.Name)

}

func main() {}
