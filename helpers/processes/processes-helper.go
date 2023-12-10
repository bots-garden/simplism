package processesHelper

import (
	"fmt"
	"os"
	"os/exec"
	configHelper "simplism/helpers/config"
	stringHelper "simplism/helpers/stringHelper"
	simplismTypes "simplism/types"
)

// getExecutablePath returns the path of the executable file for the given program name.
//
// It takes a string parameter `progName` which represents the name of the program.
// It returns a string which represents the path of the executable file.
func GetExecutablePath(progName string) string {
	executablePath, err := exec.LookPath(progName)
	if err != nil {
		fmt.Println("😡 Error finding executable:", err)
		os.Exit(1)
	}
	return executablePath
}

func SpawnSimplismProcess(wasmArguments simplismTypes.WasmArguments) {

	wasmArguments = configHelper.ApplyDefaultValuesIfMissing(wasmArguments)

	spawnArgs := []string{
		"",
		"listen",
		wasmArguments.FilePath,
		wasmArguments.FunctionName,
		"--wasm-url", wasmArguments.URL,
		"--wasm-url-auth-header", wasmArguments.WasmURLAuthHeader,
		"--http-port", wasmArguments.HTTPPort,
		"--log-level", wasmArguments.LogLevel,
		"--allow-hosts", wasmArguments.AllowHosts,
		"--allow-paths", wasmArguments.AllowPaths,
		"--env", wasmArguments.EnvVars,
		"--config", wasmArguments.Config,
		"--wasi", stringHelper.GetTheStringValueOf(wasmArguments.Wasi),
		"--input", wasmArguments.Input,
		"--wasm-url-auth-header", wasmArguments.WasmURLAuthHeader,
		"--cert-file", wasmArguments.CertFile,
		"--key-file", wasmArguments.KeyFile,
		"--admin-reload-token", wasmArguments.AdminReloadToken,
		"--admin-discovery-token", wasmArguments.AdminDiscoveryToken,
		"--admin-spawn-token", wasmArguments.AdminSpawnToken,
		//"--spawn-mode", "true",
		"--service-discovery", stringHelper.GetTheStringValueOf(wasmArguments.ServiceDiscovery),
		"--discovery-endpoint", wasmArguments.DiscoveryEndpoint,
		"--service-name", wasmArguments.ServiceName,
		"--information", wasmArguments.Information,
	}

	cmd := &exec.Cmd{
		Path:   GetExecutablePath("simplism"),
		Args:   spawnArgs,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err := cmd.Start()
	if err != nil {
		fmt.Println("😡 Error when spawning a new simplism process:", wasmArguments.FilePath, err)
		os.Exit(1) // exit with an error 🤔
		// TODO: return something instead of exiting
	}
}

func KillSimplismProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("😡 Error finding process", err)
		return err
	}

	err = process.Kill()
	if err != nil {
		fmt.Println("😡 Error killing process", err)
		return err
	}

	return nil

}