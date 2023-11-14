package main

import (
	"net/http"

	"github.com/vmware-labs/wasm-workers-server/kits/go/worker"
)

func main() {
	worker.ServeFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-generated-by", "wasm-workers-server")
		w.Write([]byte("Hello Wasm Workers Server!"))

	})
}
