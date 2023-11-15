// Package generators provides the generators for simplism.
package generators

import (
	//"embed" used for the template files
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"strings"
)

//go:embed go-template.txt
var goTemplate []byte

//go:embed go-module.txt
var goModule []byte

//go:embed go-readme.txt
var goReadMe []byte

// Generate generates a project in the specified language, with the given project name and path.
//
// Parameters:
// - language: The programming language of the project. Possible values are "golang", "ruslang", or "javascript".
// - projectName: The name of the project.
// - projectPath: The path where the project will be generated.
//
// Return type: None
func Generate(language string, projectName string, projectPath string) {
	//fmt.Println("🚧 create", language, "project, named:", projectName, "on:", projectPath)

	if projectName == "" {
		fmt.Println("😡 project name cannot be empty")
		os.Exit(1)
	}
	if projectPath == "" {
		fmt.Println("😡 project path cannot be empty")
		os.Exit(1)
	}
	if language == "" {
		fmt.Println("😡 project language cannot be empty")
		os.Exit(1)
	}
	// test the value of language, the possible values are golang, ruslang or javascript, then use switch case
	switch language {
	case "golang", "go":
		/*
			./simplism generate golang yo --path=my-projects
		*/
		// Code for generating a Go project
		fmt.Println("🚧 Generating Go project...")
		// Create a directory with projectPath + projectName
		err := os.MkdirAll(projectPath+"/"+projectName, os.ModePerm)
		if err != nil {
			fmt.Println("😡 failed to create directory:", err)
			os.Exit(1)
		}
		// Generate a file main.go into the path with string(goTemplate) as content
		err = os.WriteFile(projectPath+"/"+projectName+"/main.go", []byte(goTemplate), 0644)
		if err != nil {
			fmt.Println("😡 failed to write file:", err)
			os.Exit(1)
		}
		// Generate a file go.mod into the path with string(goModule) as content and replace <name> with projectName and <version> with the version of the go compiler
		var strGoModule = strings.Replace(string(goModule), "<name>", projectName, 1)
		strGoModule = strings.Replace(strGoModule, "<version>", runtime.Version(), 1)

		err = os.WriteFile(projectPath+"/"+projectName+"/go.mod", []byte(strGoModule), 0644)
		if err != nil {
			fmt.Println("😡 failed to write file:", err)
			os.Exit(1)
		}

		var strGoReadMe = strings.Replace(string(goReadMe), "<name>", projectName, 3)
		err = os.WriteFile(projectPath+"/"+projectName+"/README.md", []byte(strGoReadMe), 0644)
		if err != nil {
			fmt.Println("😡 failed to write file:", err)
			os.Exit(1)
		}

	case "ruslang", "rust":
		// Code for generating a Ruslang project
		fmt.Println("Generating Ruslang project...")
	case "javascript", "js":
		// Code for generating a JavaScript project
		fmt.Println("Generating JavaScript project...")
	default:
		fmt.Println("Invalid language specified.")
	}

}
