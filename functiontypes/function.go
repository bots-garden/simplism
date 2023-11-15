package functiontypes

import (
	"encoding/base64"
	"encoding/json"
)

// Argument of the remote function
type Argument struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Method string              `json:"method"`
	URI    string              `json:"uri"`
}

// ToJSONString convert the argument to a JSON string
func (m *Argument) ToJSONString() string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "error" // 🤔
	}
	return string(jsonBytes)
}

// ToJSONBuffer convert the argument to a JSON buffer
func (m *Argument) ToJSONBuffer() []byte {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return []byte("error") // 🤔
	}
	return jsonBytes
}

// ToEncodedJSONString convert the argument to a base64 encoded JSON string
func (m *Argument) ToEncodedJSONString() string {

	return base64.StdEncoding.EncodeToString([]byte(m.ToJSONString()))
}

// ReturnValue of the remote function
type ReturnValue struct {
	Body   string              `json:"body"`
	Header map[string][]string `json:"header"`
	Code   int                 `json:"code"`
}
