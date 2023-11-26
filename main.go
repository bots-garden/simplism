package main

// create a simple http server
import (
	"flag"
	"fmt"
	"os"
	"simplism/cmds"
)

// main is the entry point of the Go program.
//
// It parses the command line flags and checks if a command is provided.
// If a command is provided, it calls the corresponding command function.
// Otherwise, it prints an error message and exits.
func main() {

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("😡 invalid command")
		os.Exit(0)
	}

	command := flag.Args()[0]

	errCmd := cmds.Parse(command, flag.Args()[1:])

	if errCmd != nil {
		fmt.Println("😡", errCmd)
		os.Exit(1)
	}
}
