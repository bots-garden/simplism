package httpHelper

import (
	"net/http"
	"strings"
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

// GetHeaderFromString returns the header name and value from a given header string.
//
// The parameter headerNameAndValue is a string that contains the header name and value separated by "=".
// The function splits the headerNameAndValue string using "=" as the separator and stores the result in the splitHeader variable.
// The header name is extracted from the first element of the splitHeader array and stored in the headerName variable.
// The header value is obtained by joining all the elements of the splitHeader array, except the first one, using an empty string as the separator.
// The function then returns the headerName and headerValue as a tuple.
//
// The function returns two string values: headerName and headerValue.
func GetHeaderFromString(headerNameAndValue string) (string, string) {
	splitHeader := strings.Split(headerNameAndValue, "=")
	headerName := splitHeader[0]
	// join all item of splitAuthHeader with "" except the first one
	headerValue := strings.Join(splitHeader[1:], "")
	return headerName, headerValue
}
