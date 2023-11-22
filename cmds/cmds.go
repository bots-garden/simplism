package cmds

import (
	//"embed" used for the version.txt file
	_ "embed"
	"flag"
	"fmt"
	"os"
	"simplism/generators"
	"simplism/server"

	"gopkg.in/yaml.v3"
)

//go:embed version.txt
var version []byte

// readYamlFile reads a YAML file and returns a map of server.WasmArguments and an error.
//
// It takes a string parameter called yamlFilePath, which represents the path to the YAML file.
// The function returns a map[string]server.WasmArguments, which is a map of server.WasmArguments objects, and an error.
func readYamlFile(yamlFilePath string) (map[string]server.WasmArguments, error) {

	yamlFile, err := os.ReadFile(yamlFilePath)

	if err != nil {
		return nil, err
	}

	data := make(map[string]server.WasmArguments)

	err = yaml.Unmarshal(yamlFile, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func applyDefaultValuesIfMissing(wasmArguments server.WasmArguments) server.WasmArguments {
	// default values:
	if wasmArguments.AllowHosts == "" {
		wasmArguments.AllowHosts = `["*"]`
	}

	if wasmArguments.AllowPaths == "" {
		wasmArguments.AllowPaths = "{}"
	}

	if wasmArguments.Config == "" {
		wasmArguments.Config = "{}"
	}

	if wasmArguments.HTTPPort == "" {
		wasmArguments.HTTPPort = "8080"
	}
	return wasmArguments

}

// Parse parses the command and arguments to perform a specific action.
//
// The function takes in a command string and an array of arguments.
// The command string specifies the action to be performed.
// The args array contains additional arguments for the command.
//
// The function returns an error if there is an issue during parsing.
func Parse(command string, args []string) error {

	switch command {

	case "config":
		configFilepath := flag.Args()[1] // path of the config file

		wasmArgumentsMap, err := readYamlFile(configFilepath)
		if err != nil {
			fmt.Println("🔴 reading the yaml config file:", err)
			os.Exit(1)
		}

		if len(flag.Args()) <= 2 {
			fmt.Println("🔴 you must provide a configuration key")
			os.Exit(1)

		} else {
			configKey := flag.Args()[2]

			// Start the server with the specified wasm plugin in the config
			wasmArguments := wasmArgumentsMap[configKey]
			wasmArguments = applyDefaultValuesIfMissing(wasmArguments)
			server.Listen(wasmArguments)
		}

		//os.Exit(0)
		return nil

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

		certFile := flagSet.String("cert-file", "", "Certificate file")
		keyFile := flagSet.String("key-file", "", "Key file")

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
			CertFile:        *certFile,
			KeyFile:         *keyFile,
		})
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

	// TODO: add help, about
	default:
		return fmt.Errorf("🔴 invalid command")
	}
}
