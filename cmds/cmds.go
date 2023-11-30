package cmds

import (
	//"embed" used for the version.txt file

	_ "embed"
	"flag"
	"fmt"
	"simplism/generators"
)

//go:embed version.txt
var version []byte

//go:embed about.txt
var about []byte

// Parse parses the command and arguments to perform a specific action.
//
// The function takes in a command string and an array of arguments.
// The command string specifies the action to be performed.
// The args array contains additional arguments for the command.
//
// The function returns an error if there is an issue during parsing.
func Parse(command string, args []string) error {

	switch command {

	case "flock":
		configFilepath := flag.Args()[1] // path of the config file

		startFlockMode(configFilepath)

		return nil

	case "config":
		configFilepath := flag.Args()[1] // path of the config file

		startConfigMode(configFilepath)

		//os.Exit(0)
		return nil

	case "listen":
		wasmFilePath := flag.Args()[1]     // path of the wasm file
		wasmFunctionName := flag.Args()[2] // function name

		flagSet := flag.NewFlagSet("listen", flag.ExitOnError)

		startListening(wasmFilePath, wasmFunctionName, flagSet, args)

		return nil

	case "version":
		fmt.Println(string(version))
		//os.Exit(0)
		return nil

	case "generate":
		/*
			./simplism generate golang hello projects
		*/
		language := flag.Args()[1]    // language of the project
		projectName := flag.Args()[2] // name of the project
		projectPath := flag.Args()[3] // path of the project

		generators.Generate(language, projectName, projectPath)

		return nil
	
	case "about":
		fmt.Println(string(about))
		return nil

	// TODO: add help
	default:
		return fmt.Errorf("😡 invalid command")
	}
}
