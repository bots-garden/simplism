package function

import (
	"encoding/base64"
	"encoding/json"
)

type MainArgument struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Method string              `json:"method"`
	Uri    string              `json:"uri"`
}

func (m *MainArgument) ToJSONString() string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "error" // 🤔
	}
	return string(jsonBytes)
}

func (m *MainArgument) ToJSONBuffer() []byte {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return []byte("error") // 🤔
	}
	return jsonBytes
}

func (m *MainArgument) ToEncodedJSONString() string {

	return base64.StdEncoding.EncodeToString([]byte(m.ToJSONString()))
}

type ResponseResult struct {
}

type ReturnValue struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Code   int                 `json:"code"`
}
