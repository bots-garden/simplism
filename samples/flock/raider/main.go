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

// RaiderMessage structure (to send data to the baseStar)
type RaiderMessage struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// BaseStarMessage structure (data received from the baseStar)
type BaseStarMessage struct {
	Text string `json:"text"`
	Cmd  string `json:"cmd"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// sendResponse sends a response containing a RaiderMessage.
//
// The function takes a RaiderMessage as a parameter and constructs a ResponseData struct with the Body, Header, and Code fields set.
// It then converts the ResponseData struct to a JSON string using the `json.Marshal` function.
// The resulting JSON string is copied to host memory using the `pdk.AllocateBytes` function,
// and the memory is outputted using the `pdk.OutputMemory` function.
func sendResponse(raiderMessage RaiderMessage) {
	response := ResponseData{
		Body:   toJsonString(raiderMessage),
		Header: map[string][]string{"Content-Type": {"application/json; charset=utf-8"}},
		Code:   200,
	}
	// response to Json string
	jsonResponse, _ := json.Marshal(response)

	// copy output to host memory
	mem := pdk.AllocateBytes(jsonResponse)
	pdk.OutputMemory(mem)
}

// getRequestData reads the function argument from memory and unmarshals it into a RequestData struct.
//
// No parameters.
// Returns a RequestData object.
func getRequestData() RequestData {
	// read function argument from the memory
	requestDataBytes := pdk.Input()

	var requestData RequestData // this is the data request
	json.Unmarshal(requestDataBytes, &requestData)
	return requestData
}

// toJsonString converts the given data to a JSON string.
//
// It takes in a single parameter, `data`, which is an interface{} type representing the data to be converted.
// It returns a string, which is the JSON representation of the data.
func toJsonString(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

// sendMessageToBaseStar sends a raiderMessage to the specified baseStarURL and returns the response as a BaseStarMessage.
//
// Parameters:
// - raiderMessage: The raiderMessage to send.
// - baseStarURL: The URL of the base star to send the message to.
//
// Return Type:
// - BaseStarMessage: The response from the base star.
func sendMessageToBaseStar(raiderMessage RaiderMessage, baseStarURL string) BaseStarMessage {
	// send request to the basestar
	req := pdk.NewHTTPRequest("POST", baseStarURL)
	req.SetHeader("Content-Type", "application/json")
	req.SetBody([]byte(toJsonString(raiderMessage)))
	res := req.Send() // send message to basestar

	// display response from the basestar
	var baseStarMessage BaseStarMessage
	json.Unmarshal(res.Body(), &baseStarMessage)

	return baseStarMessage
}

// readConfig retrieves the configuration values from pdk.GetConfig for raider-name, basestar-url, raider-url,
// y-start, and x-start. It converts x-start and y-start to integers using strconv.Atoi.
//
// It returns the configuration values as strings for raiderName, baseStarURL, raiderURL and as integers for x and y.
func readConfig() (string, string, string, int, int) {
	raiderName, _ := pdk.GetConfig("raider-name")
	baseStarURL, _ := pdk.GetConfig("basestar-url")
	raiderURL, _ := pdk.GetConfig("raider-url")

	yStart, _ := pdk.GetConfig("y-start")
	xStart, _ := pdk.GetConfig("x-start")

	x, _ := strconv.Atoi(xStart)
	y, _ := strconv.Atoi(yStart)

	return raiderName, baseStarURL, raiderURL, x, y
}

// getCoordinatesFromMemory retrieves the coordinates from memory.
//
// It takes two parameters: xIfZero (an integer) and yIfZero (an integer).
// It returns two integers representing the x and y coordinates.
func getCoordinatesFromMemory(xIfZero int, yIfZero int) (int, int) {
	var x, y int
	x = pdk.GetVarInt("x")
	y = pdk.GetVarInt("y")
	if x == 0 {
		x = xIfZero
	}
	if y == 0 {
		y = yIfZero
	}
	return x, y
}

// setCoordinatesToMemory sets the given x and y coordinates to the memory.
//
// Parameters:
// - x: an integer representing the x coordinate.
// - y: an integer representing the y coordinate.
func setCoordinatesToMemory(x int, y int) {
	pdk.SetVarInt("x", x)
	pdk.SetVarInt("y", y)
}

/*
This code defines a function called handle with no parameters.

- Inside the function, it calls the readConfig function to retrieve some configuration values and assigns them to variables.
- Then, it calls the getCoordinatesFromMemory function to retrieve the current coordinates from memory, or use default values if they are not set.
- Next, it gets the request data and checks the URI.
  - If the URI is "/move", it creates a RaiderMessage struct with the retrieved values and sends the message to a baseStar using the sendMessageToBaseStar function.
    - It then logs the response from the base star and if the response command is "move",
    - it updates the coordinates in memory and in the RaiderMessage struct.
    - Finally, it sends the RaiderMessage as a response.
  - If the URI is not "/move", it creates a RaiderMessage struct and sends it as a response.

The RaiderMessage struct contains the raider name, URL, and coordinates.
The sendMessageToBaseStar function sends the RaiderMessage to a baseStar and returns the response as a BaseStarMessage struct.
The getCoordinatesFromMemory function retrieves the coordinates from memory and returns them.
The setCoordinatesToMemory function sets the given coordinates to memory.

*/
//export handle
func handle() {

	// get config information
	raiderName, baseStarURL, raiderURL, xStart, yStart := readConfig()

	// get coordinates from memory, if they are not set, use the start values
	x, y := getCoordinatesFromMemory(xStart, yStart)

	requestData := getRequestData()
	switch requestData.URI {

	case "/move":

		raiderMessage := RaiderMessage{
			Name: raiderName,
			Url:  raiderURL,
			X:    x,
			Y:    y,
		}

		// send request to the basestar (ask for coordinates)
		baseStarMessage := sendMessageToBaseStar(raiderMessage, baseStarURL+"/need-coordinates")

		// display response from the basestar
		pdk.Log(
			pdk.LogInfo,
			"📩 from basestar: "+baseStarMessage.Text+" "+baseStarMessage.Cmd+" to: "+strconv.Itoa(baseStarMessage.X)+":"+strconv.Itoa(baseStarMessage.Y),
		)

		if baseStarMessage.Cmd == "move" {
			// new coordinates
			setCoordinatesToMemory(baseStarMessage.X, baseStarMessage.Y)

			raiderMessage.X = baseStarMessage.X
			raiderMessage.Y = baseStarMessage.Y
		}

		sendResponse(raiderMessage)

	default:
		raiderMessage := RaiderMessage{
			Name: raiderName,
			Url:  raiderURL,
			X:    x,
			Y:    y,
		}
		sendResponse(raiderMessage)
	}

}

func main() {}
