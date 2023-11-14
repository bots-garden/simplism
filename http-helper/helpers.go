package httpHelper

import (
	"encoding/json"
	"net/http"
	"strings"
)

func HeaderToJSONString(header http.Header) string {
	jsonBytes, err := json.Marshal(header)
	if err != nil {
		// handle error
	}

	return string(jsonBytes)
}

func IsContentTypeJson(header http.Header) bool {
	return strings.HasPrefix(header.Get("Content-Type"), "application/json") || strings.HasPrefix(header.Get("content-type"), "application/json")
}
func IsContentTypeText(header http.Header) bool {
	return strings.HasPrefix(header.Get("Content-Type"), "text/plain") || strings.HasPrefix(header.Get("content-type"), "text/plain")
}

func IfContentTypeIsJson(header http.Header, callback func()) bool {
	if IsContentTypeJson(header) {
		callback()
		return true
	}
	return false
}

func IfContentTypeIsText(header http.Header, callback func()) bool {
	if IsContentTypeText(header) {
		callback()
		return true
	}
	return false
}

func GetBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}
