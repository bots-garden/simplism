// Package wasmhelper contains helper functions for the wasm runtime
package wasmhelper

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

// WasmPlugin type
type WasmPlugin struct {
    ExtismPlugin *extism.Plugin
    Protection   sync.Mutex
}

var (
    wasmPlugins     = make(map[string]*WasmPlugin)
    counter         = 0
    poolSize        = 4
    prefixPluginKey = "plugin"
)

// GetLevel returns the corresponding extism.LogLevel based on the given log level string.
//
// It takes in a logLevel string and returns an extism.LogLevel.
func GetLevel(logLevel string) extism.LogLevel {
    level := extism.Off
    switch logLevel {
    case "error":
        level = extism.Error
    case "warn":
        level = extism.Warn
    case "info":
        level = extism.Info
    case "debug":
        level = extism.Debug
    case "trace":
        level = extism.Trace
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
        LogLevel:     &logLevel,
    }

    manifest := extism.Manifest{
        Wasm: []extism.Wasm{
            extism.WasmFile{
                Path: wasmFilePath},
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
    pluginInst, err := extism.NewPlugin(ctx, manifest, config, nil) // new
    if err != nil {
        log.Println("😡 Error when creating the wasm plugin instance", err)
        os.Exit(1)
    }

    wasmPlugin := WasmPlugin{
        ExtismPlugin: pluginInst,
    }
    return &wasmPlugin
}

// GeneratePluginsPool generates a pool of plugins.
//
// It takes the following parameters:
// - ctx: the context.Context object for cancellation and timeouts.
// - config: the PluginConfig object containing the configuration for the plugins.
// - manifest: the Manifest object containing the manifest of the plugins.
//
// This function does not return any value.
func GeneratePluginsPool(ctx context.Context, config extism.PluginConfig, manifest extism.Manifest) {
    for i := 0; i <= poolSize; i++ {
        wasmPlugin := getPluginInstance(ctx, config, manifest)
        key := prefixPluginKey + strconv.Itoa(i)
        storePlugin(key, wasmPlugin)
    }
}

// CallWasmFunction executes a WebAssembly function.
//
// CallWasmFunction takes in the name of the WebAssembly function as a
// string and the parameters to be passed to the function as a byte slice.
// It returns the output of the function as a byte slice and any error
// encountered during the execution of the function.
func CallWasmFunction(wasmFunctionName string, params []byte) ([]byte, error) {
    key := prefixPluginKey + strconv.Itoa(counter)
    wasmPlugin := getPlugin(key)
    counter ++
    if counter == poolSize {
        counter = 0
    }
    wasmPlugin.Protection.Lock()

    defer wasmPlugin.Protection.Unlock()

    _, out, err := wasmPlugin.ExtismPlugin.Call(wasmFunctionName, params)

    if err != nil {
        return nil, err
    } else {
        return out, nil
    }

}
