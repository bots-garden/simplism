// main package
package main

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"

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

// RaiderMessage structure (data sent by the Raider)
type RaiderMessage struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// BaseStarMessage structure (data sent to the raider by the baseStar)
type BaseStarMessage struct {
	Text string `json:"text"`
	Cmd  string `json:"cmd"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type Message struct {
	Text string `json:"text"`
}

// displayCounter is a Go function that displays the value of the counter.
//
// It does not have any parameters.
// It does not have a return type.
func displayCounter() {
	count := pdk.GetVarInt("count")
	pdk.Log(pdk.LogInfo, "⭐️ counter: "+strconv.Itoa(count))
}

// countAccess increments the value of the "count" variable by 1.
//
// This function does not take any parameters.
// It does not return any values.
func countAccess() {
	var count int
	count = pdk.GetVarInt("count")
	count = count + 1
	pdk.SetVarInt("count", count)
}

// getRequestData retrieves the request data from memory and parses it into a RequestData struct.
//
// This function does not take any parameters.
// It returns a RequestData struct.
func getRequestData() RequestData {
	// read function argument from the memory
	requestDataBytes := pdk.Input()

	var requestData RequestData // this is the data request
	json.Unmarshal(requestDataBytes, &requestData)
	return requestData
}

// getRaiderMessage unmarshalls the requestData.Body and returns the raider message.
//
// It takes a requestData of type RequestData as a parameter.
// It returns a RaiderMessage.
func getRaiderMessage(requestData RequestData) RaiderMessage {
	// unmarshall and get the value of the raider message
	var raiderMessage RaiderMessage
	json.Unmarshal([]byte(requestData.Body), &raiderMessage)
	return raiderMessage
}

// sendResponse sends a response to the raider.
//
// It takes a message parameter of type string.
// It does not return anything.
func sendResponse(message string) {
	// send response to the raider
	response := ResponseData{
		Body:   message,
		Header: map[string][]string{"Content-Type": {"application/json; charset=utf-8"}},
		Code:   200,
	}
	// response to Json string
	jsonResponse, _ := json.Marshal(response)

	// copy output to host memory
	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)
}

// toJsonString converts the given data into a JSON string.
//
// The data parameter is the input data to be converted into a JSON string.
// It can be of any type.
//
// The return type is a string that represents the JSON string of the input data.
func toJsonString(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

// getRaidersData returns a map of raider names to their corresponding RaiderMessage.
//
// This function does the following:
// - Retrieves the list of all raiders from the config.
// - Loops through the list of raiders.
// - Retrieves the raider message from the variables using json.Unmarshal.
// - Adds the raider name and message to the raidersMap.
//
// Returns:
// - A map[string]RaiderMessage: a map of raider names to their corresponding RaiderMessage.
func getRaidersData() map[string]RaiderMessage {
	// get the list of all raiders
	var raidersMap = make(map[string]RaiderMessage)

	configRaiders, _ := pdk.GetConfig("raiders")
	raiders := strings.Split(configRaiders, ",")

	// loop through the list of raiders
	for _, raiderName := range raiders {
		var raiderMessage RaiderMessage
		json.Unmarshal(pdk.GetVar(raiderName), &raiderMessage)
		raidersMap[raiderName] = raiderMessage
	}
	return raidersMap
}

// getXYMax returns the maximum values for x and y.
//
// This function does not take any parameters.
// It returns two integers, representing the maximum values for x and y.
func getXYMax() (int, int) {
	xMax, _ := pdk.GetConfig("x-max")
	yMax, _ := pdk.GetConfig("y-max")

	// cast xMax and yMax to integers
	xMaxInt, _ := strconv.Atoi(xMax)
	yMaxInt, _ := strconv.Atoi(yMax)

	return xMaxInt, yMaxInt
}

// return random integer between -1 and 1
func randomInt() int {
	return rand.Intn(3) - 1
}

/*
This code performs the following steps:

- Calls the countAccess() function.
- Calls the displayCounter() function.
- Retrieves the request data using the getRequestData() function.
- Uses a switch statement to handle different URIs:
  - If the URI is "/raiders",
    - it sends a response by converting the raiders data to JSON.
  - If the URI is "/need-coordinates",
    - it creates a BaseStarMessage struct,
	- updates the coordinates,
	- stores the raider message in memory,
	- and sends the updated BaseStarMessage as the response.
  - For any other URI,
	- it gets the base star name from the configuration,
	- creates a Message struct,
	- and sends a response with a greeting from the base star.
*/
//export handle
func handle() {
	countAccess()
	displayCounter()

	requestData := getRequestData()
	switch requestData.URI {
	case "/raiders":
		sendResponse(toJsonString(getRaidersData()))

	case "/need-coordinates":
		raiderMessage := getRaiderMessage(requestData)

		xMax, yMax := getXYMax() // borders

		raiderMessage.X = raiderMessage.X + randomInt()
		raiderMessage.Y = raiderMessage.Y + randomInt()

		if raiderMessage.X > xMax {
			raiderMessage.X = raiderMessage.X - 1
		}
		if raiderMessage.X < 0 {
			raiderMessage.X = raiderMessage.X + 1
		}
		if raiderMessage.Y > yMax {
			raiderMessage.Y = raiderMessage.Y -1
		}
		if raiderMessage.Y < 0 {
			raiderMessage.Y = raiderMessage.Y +1
		}

		baseStarMessage := BaseStarMessage{
			Cmd:  "move",
			Text: "🤗 Hello ⭐️" + raiderMessage.Name,
			X:    raiderMessage.X,
			Y:    raiderMessage.Y,
		}

		// store the data in the memory
		pdk.SetVar(raiderMessage.Name, []byte(toJsonString(raiderMessage)))

		sendResponse(toJsonString(baseStarMessage))

	default:
		baseStarName, _ := pdk.GetConfig("basestar-name")
		sendResponse(toJsonString(Message{Text: "👋 Hello from " + baseStarName}))
	}

}

func main() {}
