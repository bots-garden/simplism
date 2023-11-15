package httpHelper

import (
	"net/http"
)

// GetBody returns the body of an HTTP request as a byte slice.
//
// It takes a pointer to an http.Request as a parameter.
// It returns a byte slice.
func GetBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}
