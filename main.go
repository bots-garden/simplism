package main

// create a simple http server
import (
	"flag"
	"fmt"
	"os"
	"simplism/cmds"
)

/*
go run main.go \
listen ./samples/hello-plugin/simple.wasm \
say_hello \
  --http-port 8080 \
  --log-level info
*/


func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("🔴 invalid command")
		os.Exit(0)
	}

	command := flag.Args()[0]

	errCmd := cmds.Parse(command, flag.Args()[1:])

	if errCmd != nil {
		fmt.Println(errCmd)
		os.Exit(1)
	}
}
