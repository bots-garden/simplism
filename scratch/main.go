// main package
package main

import (
	"github.com/extism/go-pdk"
)

//export handle
func handle() {

	mem := pdk.AllocateBytes([]byte(
		`{"body":"🖖 Live long and prosper 🤗 | simplism v0.1.3","header":{"Content-Type":["text/plain; charset=utf-8"]},"code":200}`,
	))
	pdk.OutputMemory(mem)

}

func main() {}
