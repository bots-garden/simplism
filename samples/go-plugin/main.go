package main

import (
	"strconv"

	"github.com/extism/go-pdk"
	"github.com/valyala/fastjson"
)

var parser = fastjson.Parser{}

func GetJsonFromBytes(buffer []byte) (*fastjson.Value, error) {
	return parser.ParseBytes(buffer)
}

func GetJsonFromString(jsonString string) (*fastjson.Value, error) {
	return parser.Parse(jsonString)
}

func TextPlainResponse(bodyContent string, httpCode int) string {
	return `{"code":` + strconv.Itoa(httpCode) + `, "body":"` + bodyContent + `","header":{"Content-Type":["text/plain; charset=utf-8"]}}`
}

//export say_hello
func say_hello() {
	// read function argument from the memory
	input := pdk.Input()

	jsonArgument, _ := GetJsonFromBytes(input)

	method := string(jsonArgument.GetStringBytes("method"))
	//uri := string(jsonArgument.GetStringBytes("uri"))

	//header := jsonArgument.GetObject("header")
	//contentType := string(header.Get("Content-Type").GetArray()[0].GetStringBytes())

	var message = ""
	var body = "{}"
	if method == "POST" {

		body = string(jsonArgument.GetStringBytes("body"))

		//pdk.Log(pdk.LogInfo, "📙 body: "+body)

		jsonBody, _ := GetJsonFromString(body)

		firstName := jsonBody.GetStringBytes("firstName")
		lastName := jsonBody.GetStringBytes("lastName")

		message = "👋 Hello " + string(firstName) + " " + string(lastName)

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
	}
	*/

	response := TextPlainResponse(message, 200)

	// copy output to host memory
	mem := pdk.AllocateString(response)
	pdk.OutputMemory(mem)

}

func main() {}

/*
👋 don't forget the encoded string like with slingshot
*/
