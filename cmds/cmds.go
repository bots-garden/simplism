package cmds

import (
	//"embed" used for the version.txt file
	_ "embed"
	"flag"
	"fmt"
	"simplism/generators"
	"simplism/server"
)

//go:embed version.txt
var version []byte

// Parse parses the command and arguments to perform a specific action.
//
// The function takes in a command string and an array of arguments.
// The command string specifies the action to be performed.
// The args array contains additional arguments for the command.
//
// The function returns an error if there is an issue during parsing.
func Parse(command string, args []string) error {

	switch command {

	case "listen":

		wasmFilePath := flag.Args()[1]     // path of the wasm file
		wasmFunctionName := flag.Args()[2] // function name

		flagSet := flag.NewFlagSet("listen", flag.ExitOnError)

		httpPort := flagSet.String("http-port", "8080", "http port")

		input := flagSet.String("input", "", "Argument of the function")
		logLevel := flagSet.String("log-level", "", "Log level to print message")
		allowHosts := flagSet.String("allow-hosts", `["*"]`, "Hosts for HTTP request (json array)")
		allowPaths := flagSet.String("allow-paths", "{}", "Allowed paths to write and read files (json string)")
		config := flagSet.String("config", "{}", "Configuration data (json string)")
		wasi := flagSet.Bool("wasi", true, "")

		wasmURL := flagSet.String("wasm-url", "", "Url to download the wasm file")
		authHeaderName := flagSet.String("auth-header-name", "", "Authentication header name, ex: PRIVATE-TOKEN")
		authHeaderValue := flagSet.String("auth-header-value", "", "Value of the authentication header, ex: IlovePandas")

		flagSet.Parse(args[2:])

		server.Listen(server.WasmArguments{
			FilePath:        wasmFilePath,
			FunctionName:    wasmFunctionName,
			HTTPPort:        *httpPort,
			Input:           *input,
			LogLevel:        *logLevel,
			AllowHosts:      *allowHosts,
			AllowPaths:      *allowPaths,
			Config:          *config,
			Wasi:            *wasi,
			URL:             *wasmURL,
			AuthHeaderName:  *authHeaderName,
			AuthHeaderValue: *authHeaderValue,
		})
		return nil

	case "version":
		fmt.Println(string(version))
		//os.Exit(0)
		return nil 
	
	case "generate":
		/*
			./simplism generate golang hello --path=./projects
		*/
		language := flag.Args()[1]     // language of the project
		projectName := flag.Args()[2] // name of the project

		flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
		projectPath := flagSet.String("path", "", "Project path")

		flagSet.Parse(args[2:])

		generators.Generate(language, projectName, *projectPath)


		return nil

	// TODO: add help, about, generators, ...
	default:
		return fmt.Errorf("🔴 invalid command")
	}
}
