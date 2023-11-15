package wasmHelper

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

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

func GetConfigAndManifest(wasmFilePath string, hosts []string, paths map[string]string, manifestConfig map[string]string, logLevel extism.LogLevel) (extism.PluginConfig, extism.Manifest) {
	//logLevel := GetLevel("info") // tmp

	//fmt.Println(logLevel)

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

/*
func GetCounter() int {
	return counter
}

func GetPoolSize() int {
	return poolSize
}
*/

func storePlugin(key string, wasmPlugin *WasmPlugin) {
	wasmPlugins[key] = wasmPlugin
}

func getPlugin(key string) *WasmPlugin {
	return wasmPlugins[key]
}

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

func GeneratePluginsPool(ctx context.Context, config extism.PluginConfig, manifest extism.Manifest) {
	for i := 0; i <= poolSize; i++ {
		wasmPlugin := getPluginInstance(ctx, config, manifest)
		key := prefixPluginKey + strconv.Itoa(i)
		storePlugin(key, wasmPlugin)
	}
}

func CallWasmFunction(wasmFunctionName string, params []byte) ([]byte, error) {
	key := prefixPluginKey + strconv.Itoa(counter)
	wasmPlugin := getPlugin(key)
	counter += 1
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
