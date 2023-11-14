package check

import "log"

// Are there enough arguments?
func IfThereAreEnoughArgs(args []string) {
	switch numberOfArgs := len(args); {
	case numberOfArgs == 0:
		log.Fatal("😡 no wasm file path")
	case numberOfArgs == 1:
		log.Fatal("😡 no wasm function name")
	case numberOfArgs == 2:
		log.Fatal("😡 no http port")
	}
}
