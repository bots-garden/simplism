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

//----------------------
// Golang
//----------------------

//go:embed go-template.txt
var goTemplate []byte

//go:embed go-module.txt
var goModule []byte

//go:embed go-readme.txt
var goReadMe []byte

//go:embed go-build.txt
var goBuild []byte

//go:embed go-run.txt
var goRun []byte

//go:embed go-query.txt
var goQuery []byte

//----------------------
// Rustlang
//----------------------

//go:embed rust-template.txt
var rustTemplate []byte

//go:embed rust-cargo.txt
var rustCargo []byte

//go:embed rust-readme.txt
var rustReadMe []byte

//go:embed rust-build.txt
var rustBuild []byte

//go:embed rust-run.txt
var rustRun []byte

//go:embed rust-query.txt
var rustQuery []byte

//----------------------
// JavaScript
//----------------------

//go:embed js-template.txt
var jsTemplate []byte

//go:embed js-index-ts-template.txt
var jsIndexTsTemplate []byte

//go:embed js-build.txt
var jsBuild []byte

//go:embed js-run.txt
var jsRun []byte

//go:embed js-query.txt
var jsQuery []byte

//go:embed js-readme.txt
var jsReadMe []byte

//----------------------
// Dockerfile
//----------------------

//go:embed docker-template.txt
var dockerTemplate []byte

//go:embed docker.image.txt
var dockerImage []byte

//----------------------
// Gitpod
//----------------------

//go:embed gitpod-template.txt
var gitpodTemplate []byte

//go:embed gitpod.image.txt
var gitpodImage []byte

//go:embed simplism.version.txt
var simplismVersion []byte

// makeDirectoryStructure creates a directory with the given projectPath and projectName.
//
// Parameters:
// - projectPath: the path where the directory should be created.
// - projectName: the name of the directory to be created.
//
// Return type: none.
func makeDirectoryStructure(projectPath, projectName string, otherDirectories ...string) {

	if len(otherDirectories) > 0 {
		for _, directory := range otherDirectories {
			err := os.MkdirAll(projectPath+"/"+projectName+"/"+directory, os.ModePerm)
			if err != nil {
				fmt.Println("😡 failed to create directory:", err)
				os.Exit(1)
			}
		}
	}

	// Create a directory with projectPath + projectName
	err := os.MkdirAll(projectPath+"/"+projectName, os.ModePerm)
	if err != nil {
		fmt.Println("😡 failed to create directory:", err)
		os.Exit(1)
	}
}

// createFileFromTemplate writes a file to the specified path using the provided template.
//
// Parameters:
// - projectPath: the path to the project directory.
// - projectName: the name of the project.
// - filePath: the relative path to the file.
// - template: the content of the template file.
func createFileFromTemplate(projectPath, projectName, filePath string, template []byte) {
	err := os.WriteFile(projectPath+"/"+projectName+"/"+filePath, template, 0644)
	if err != nil {
		fmt.Println("😡 failed to write file:", err)
		os.Exit(1)
	}
}

func createBashFileFromTemplate(projectPath, projectName, filePath string, template []byte) {
	err := os.WriteFile(projectPath+"/"+projectName+"/"+filePath, template, 0777)
	if err != nil {
		fmt.Println("😡 failed to write file:", err)
		os.Exit(1)
	}
}

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

	switch language {
	case "golang", "go":
		/*
            Generating a Golang project:
            ./simplism generate golang hello my-projects/golang
		*/
		fmt.Println("🔵 Generating Go project...")

		makeDirectoryStructure(projectPath, projectName)

		// Dockerfile
		var strDockerFile = strings.Replace(string(dockerTemplate), "<name>", projectName, 2)
		strDockerFile = strings.Replace(strDockerFile, "<image>", string(dockerImage), 1)
		createFileFromTemplate(projectPath, projectName, "Dockerfile", []byte(strDockerFile))

		// .gitpod.yml
		var strGitpodYamlFile = strings.Replace(string(gitpodTemplate), "<image>", string(gitpodImage), 1)
		strGitpodYamlFile = strings.Replace(strGitpodYamlFile, "<version>", string(simplismVersion), 1)
		createFileFromTemplate(projectPath, projectName, ".gitpod.yml", []byte(strGitpodYamlFile))

		createFileFromTemplate(projectPath, projectName, "main.go", goTemplate)

		var strGoModule = strings.Replace(string(goModule), "<name>", projectName, 1)
		strGoModule = strings.Replace(strGoModule, "<version>", strings.ReplaceAll(runtime.Version(), "go", ""), 1)

		createFileFromTemplate(projectPath, projectName, "go.mod", []byte(strGoModule))

		var strGoReadMe = strings.Replace(string(goReadMe), "<name>", projectName, 6)
		createFileFromTemplate(projectPath, projectName, "README.md", []byte(strGoReadMe))

		var strGoBuild = strings.Replace(string(goBuild), "<name>", projectName, 1)
		createBashFileFromTemplate(projectPath, projectName, "build.sh", []byte(strGoBuild))

		var strGoRun = strings.Replace(string(goRun), "<name>", projectName, 1)
		createBashFileFromTemplate(projectPath, projectName, "run.sh", []byte(strGoRun))

		createBashFileFromTemplate(projectPath, projectName, "query.sh", goQuery)

		fmt.Println("🎉", "project generated in", projectPath+"/"+projectName)

	case "rustlang", "rust":
		/*
            Generating a Ruslang project:

            ./simplism generate rustlang hello my-projects/rustlang
		*/
		fmt.Println("🦀 Generating Ruslang project...")

		makeDirectoryStructure(projectPath, projectName, "src")

		// Dockerfile
		var strDockerFile = strings.Replace(string(dockerTemplate), "<name>", projectName, 2)
		strDockerFile = strings.Replace(strDockerFile, "<image>", string(dockerImage), 1)
		createFileFromTemplate(projectPath, projectName, "Dockerfile", []byte(strDockerFile))

		// .gitpod.yml
		var strGitpodYamlFile = strings.Replace(string(gitpodTemplate), "<image>", string(gitpodImage), 1)
		strGitpodYamlFile = strings.Replace(strGitpodYamlFile, "<version>", string(simplismVersion), 1)
		createFileFromTemplate(projectPath, projectName, ".gitpod.yml", []byte(strGitpodYamlFile))

		createFileFromTemplate(projectPath, projectName, "src/lib.rs", rustTemplate)

		var strRustCargo = strings.Replace(string(rustCargo), "<name>", projectName, 1)
		createFileFromTemplate(projectPath, projectName, "Cargo.toml", []byte(strRustCargo))

		var strRustReadMe = strings.Replace(string(rustReadMe), "<name>", projectName, 6)
		createFileFromTemplate(projectPath, projectName, "README.md", []byte(strRustReadMe))

		createBashFileFromTemplate(projectPath, projectName, "build.sh", rustBuild)

		var strRustRun = strings.Replace(string(rustRun), "<name>", projectName, 1)
		createBashFileFromTemplate(projectPath, projectName, "run.sh", []byte(strRustRun))

		createBashFileFromTemplate(projectPath, projectName, "query.sh", rustQuery)

		fmt.Println("🎉", "project generated in", projectPath+"/"+projectName)

	case "javascript", "js":
		fmt.Println("🟨 Generating JavaScript project...")

		makeDirectoryStructure(projectPath, projectName)

		// Dockerfile
		var strDockerFile = strings.Replace(string(dockerTemplate), "<name>", projectName, 2)
		strDockerFile = strings.Replace(strDockerFile, "<image>", string(dockerImage), 1)
		createFileFromTemplate(projectPath, projectName, "Dockerfile", []byte(strDockerFile))

		// .gitpod.yml
		var strGitpodYamlFile = strings.Replace(string(gitpodTemplate), "<image>", string(gitpodImage), 1)
		strGitpodYamlFile = strings.Replace(strGitpodYamlFile, "<version>", string(simplismVersion), 1)
		createFileFromTemplate(projectPath, projectName, ".gitpod.yml", []byte(strGitpodYamlFile))

		createFileFromTemplate(projectPath, projectName, "index.js", jsTemplate)
		createFileFromTemplate(projectPath, projectName, "index.d.ts", jsIndexTsTemplate)

		var strJsReadMe = strings.Replace(string(jsReadMe), "<name>", projectName, 6)
		createFileFromTemplate(projectPath, projectName, "README.md", []byte(strJsReadMe))

		var strJsBuild = strings.Replace(string(jsBuild), "<name>", projectName, 1)
		createBashFileFromTemplate(projectPath, projectName, "build.sh", []byte(strJsBuild))

		var strJsRun = strings.Replace(string(jsRun), "<name>", projectName, 1)
		createBashFileFromTemplate(projectPath, projectName, "run.sh", []byte(strJsRun))

		createBashFileFromTemplate(projectPath, projectName, "query.sh", jsQuery)


	default:
		fmt.Println("Invalid language specified.")
	}

}
