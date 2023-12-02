package httpHelper

import (
	"net/http"
	"os"
	simplismTypes "simplism/types"
)

// CheckDiscoveryToken checks if the provided request is authorized using the admin-discovery-token.
//
// Parameters:
// - request: the HTTP request to check.
// - wasmArgs: the Wasm arguments.
//
// Returns:
// - bool: true if the request is authorized, false otherwise.
func CheckDiscoveryToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false
	// read the header admin-discovery-token
	adminDiscoveryToken := request.Header.Get("admin-discovery-token")

	envAdminDiscoveryToken := os.Getenv("ADMIN_DISCOVERY_TOKEN")

	switch {
	// a token is awaited
	case wasmArgs.AdminDiscoveryToken != "":
		if wasmArgs.AdminDiscoveryToken == adminDiscoveryToken {
			authorised = true
		} else {
			authorised = false
		}
	// a token is awaited
	case wasmArgs.AdminDiscoveryToken == "" && envAdminDiscoveryToken != "":
		if envAdminDiscoveryToken == adminDiscoveryToken {
			authorised = true
		} else {
			authorised = false
		}
	case wasmArgs.AdminDiscoveryToken == "" && envAdminDiscoveryToken == "":
		authorised = true
	default:
		authorised = false
	}

	return authorised
}

// CheckReloadToken checks if the provided request has a valid admin reload token.
//
// It takes in the following parameters:
// - request: a pointer to an http.Request object representing the incoming request.
// - wasmArgs: a WasmArguments object representing the Wasm arguments.
//
// It returns a boolean value indicating whether the request is authorized or not.
func CheckReloadToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false

	// read the header admin-reload-token
	adminReloadToken := request.Header.Get("admin-reload-token")

	if wasmArgs.AdminReloadToken != "" {
		// token is awaited
		if wasmArgs.AdminReloadToken == adminReloadToken {
			authorised = true
		} else {
			authorised = false
		}

	} else {
		// check if the environment variable ADMIN_RELOAD_TOKEN is set
		envAdminReloadToken := os.Getenv("ADMIN_RELOAD_TOKEN")
		if envAdminReloadToken != "" {
			// token is awaited
			if envAdminReloadToken == adminReloadToken {
				authorised = true
			} else {
				authorised = false
			}
		} else {
			authorised = true
		}
	}
	return authorised
}
