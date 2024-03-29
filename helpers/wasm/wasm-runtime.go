// Package wasmhelper contains helper functions for the wasm runtime
package wasmHelper

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/go-resty/resty/v2"
	"github.com/tetratelabs/wazero"

	httpHelper "simplism/helpers/http"
	simplismTypes "simplism/types"
)

// WasmPlugin type
type WasmPlugin struct {
	ExtismPlugin *extism.Plugin
	Protection   sync.Mutex
}

var (
	wasmPlugins = make(map[string]*WasmPlugin)
	pluginKey   = "plugin"
)

// GetLevel returns the corresponding extism.LogLevel based on the given log level string.
//
// It takes in a logLevel string and returns an extism.LogLevel.
func GetLevel(logLevel string) extism.LogLevel {
	level := extism.LogLevelOff
	switch logLevel {
	case "error":
		level = extism.LogLevelError
	case "warn":
		level = extism.LogLevelWarn
	case "info":
		level = extism.LogLevelInfo
	case "debug":
		level = extism.LogLevelDebug
	case "trace":
		level = extism.LogLevelTrace
	}
	return level
}

// GetConfigAndManifest returns the plugin configuration and manifest.
//
// Parameters:
// - wasmFilePath: the path to the WebAssembly file.
// - hosts: a slice of allowed hosts.
// - paths: a map of allowed paths.
// - manifestConfig: a map of configuration values for the manifest.
// - logLevel: the log level for the plugin.
//
// Returns:
// - config: the plugin configuration.
// - manifest: the plugin manifest.
func GetConfigAndManifest(wasmFilePath string, hosts []string, paths map[string]string, manifestConfig map[string]string, logLevel extism.LogLevel) (extism.PluginConfig, extism.Manifest) {

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
		LogLevel:     logLevel,
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath,
				//Hash: "",
				//Name: "main",
			},
		},
		AllowedHosts: hosts,
		AllowedPaths: paths,
		Config:       manifestConfig,
	}

	return config, manifest
}

// storePlugin stores a WasmPlugin in the wasmPlugins map with the given key.
//
// Parameters:
// - key: the key used to store the WasmPlugin in the map.
// - wasmPlugin: the WasmPlugin to be stored.
func storePlugin(key string, wasmPlugin *WasmPlugin) {
	wasmPlugins[key] = wasmPlugin
}

// getPlugin retrieves a WasmPlugin based on the given key.
//
// Parameters:
// - key: a string representing the key used to lookup the WasmPlugin.
//
// Return:
// - *WasmPlugin: a pointer to the WasmPlugin that matches the given key.
func getPlugin(key string) *WasmPlugin {
	return wasmPlugins[key]
}

// getPluginInstance creates a new instance of WasmPlugin based on the given context, PluginConfig, and Manifest.
//
// Parameters:
//   - ctx: The context used for creating the plugin instance.
//   - config: The PluginConfig used for creating the plugin instance.
//   - manifest: The Manifest used for creating the plugin instance.
//
// Returns:
//   - *WasmPlugin: The newly created instance of WasmPlugin.
func getPluginInstance(ctx context.Context, config extism.PluginConfig, manifest extism.Manifest) *WasmPlugin {

	pluginInst, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{}) // new

	if err != nil {
		log.Println("😡 Error when creating the wasm plugin instance:", err)
		os.Exit(1)
	}
	//defer pluginInst.Close()

	wasmPlugin := WasmPlugin{
		ExtismPlugin: pluginInst,
	}
	return &wasmPlugin
}

func StartWasmPlugin(ctx context.Context, config extism.PluginConfig, manifest extism.Manifest) {
	wasmPlugin := getPluginInstance(ctx, config, manifest)
	storePlugin(pluginKey, wasmPlugin)
}

// CallWasmFunction executes a WebAssembly function.
//
// CallWasmFunction takes in the name of the WebAssembly function as a
// string and the parameters to be passed to the function as a byte slice.
// It returns the output of the function as a byte slice and any error
// encountered during the execution of the function.
func CallWasmFunction(wasmFunctionName string, params []byte) ([]byte, error) {

	key := pluginKey
	wasmPlugin := getPlugin(key)

	wasmPlugin.Protection.Lock()

	defer wasmPlugin.Protection.Unlock()

	_, out, err := wasmPlugin.ExtismPlugin.Call(wasmFunctionName, params)

	if err != nil {
		return nil, err
	} else {
		return out, nil
	}

}

func GetPlugin(index int) *WasmPlugin {
	return wasmPlugins[pluginKey+strconv.Itoa(index)]
}

// downloadWasmFile downloads a WebAssembly (Wasm) file from a given URL and saves it to the specified file path.
//
// It takes a WasmArguments struct as a parameter, which contains the necessary information for the download, such as the URL, authentication header, and file path.
// The WasmArguments struct has the following fields:
// - AuthHeaderName (string): the name of the authentication header (e.g., "PRIVATE-TOKEN")
// - AuthHeaderValue (string): the value of the authentication header (e.g., "${GITLAB_WASM_TOKEN}")
// - FilePath (string): the file path where the downloaded Wasm file will be saved
// - URL (string): the URL from which the Wasm file will be downloaded
//
// This function returns an error if there is any issue during the download process, such as a network error or an error response from the server.
// If the download is successful, it returns nil.
func DownloadWasmFile(wasmArgs simplismTypes.WasmArguments) error {
	// authenticationHeader:
	// Example: "PRIVATE-TOKEN: ${GITLAB_WASM_TOKEN}"
	client := resty.New()

	if wasmArgs.WasmURLAuthHeader != "" {
		authHeaderName, authHeaderValue := httpHelper.GetHeaderFromString(wasmArgs.WasmURLAuthHeader)
		client.SetHeader(authHeaderName, authHeaderValue)

	} else {
		// check if the environment variable WASM_URL_AUTH_HEADER is set
		wasmURLAuthHeader := os.Getenv("WASM_URL_AUTH_HEADER")
		if wasmURLAuthHeader != "" {
			authHeaderName, authHeaderValue := httpHelper.GetHeaderFromString(wasmURLAuthHeader)
			client.SetHeader(authHeaderName, authHeaderValue)

		}
	}

	resp, err := client.R().
		SetOutput(wasmArgs.FilePath).
		Get(wasmArgs.URL)

	if resp.IsError() {
		return errors.New("😡 error while downloading the wasm file, you should check the authentication token")
	}

	if err != nil {
		return err
	}
	return nil
}
