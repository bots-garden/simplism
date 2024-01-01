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

// CheckSpawnToken checks if the provided spawn token is authorized.
//
// It takes a *http.Request object and a simplismTypes.WasmArguments object as parameters.
// The function returns a boolean value indicating whether the spawn token is authorized or not.
func CheckSpawnToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false
	// read the header admin-discovery-token
	adminSpawnToken := request.Header.Get("admin-spawn-token")

	envAdminSpawnToken := os.Getenv("ADMIN_SPAWN_TOKEN")

	switch {
	// a token is awaited
	case wasmArgs.AdminSpawnToken != "":
		if wasmArgs.AdminSpawnToken == adminSpawnToken {
			authorised = true
		} else {
			authorised = false
		}
	// a token is awaited
	case wasmArgs.AdminSpawnToken == "" && envAdminSpawnToken != "":
		if envAdminSpawnToken == adminSpawnToken {
			authorised = true
		} else {
			authorised = false
		}
	case wasmArgs.AdminSpawnToken == "" && envAdminSpawnToken == "":
		authorised = true
	default:
		authorised = false
	}

	return authorised
}

// CheckStoreToken checks if the provided token is authorized to access the store.
//
// Parameters:
// - request: the HTTP request object.
// - wasmArgs: the Wasm arguments.
//
// Returns:
// - true if the token is authorized, false otherwise.
func CheckStoreToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false
	// read the header admin-store-token
	adminStoreToken := request.Header.Get("admin-store-token")

	envAdminStoreToken := os.Getenv("ADMIN_STORE_TOKEN")

	switch {
	// a token is awaited
	case wasmArgs.AdminStoreToken != "":
		if wasmArgs.AdminStoreToken == adminStoreToken {
			authorised = true
		} else {
			authorised = false
		}
	// a token is awaited
	case wasmArgs.AdminStoreToken == "" && envAdminStoreToken != "":
		if envAdminStoreToken == adminStoreToken {
			authorised = true
		} else {
			authorised = false
		}
	case wasmArgs.AdminStoreToken == "" && envAdminStoreToken == "":
		authorised = true
	default:
		authorised = false
	}

	return authorised
}

// CheckAdminRegistryToken checks if the provided request and wasm arguments have an authorized admin registry token.
//
// Parameters:
// - request: *http.Request - The HTTP request object that contains the admin registry token in the header.
// - wasmArgs: simplismTypes.WasmArguments - The Wasm arguments object that contains the admin registry token.
//
// Returns:
// - bool - True if the provided tokens are authorized, false otherwise.
func CheckAdminRegistryToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false
	// read the header admin-registry-token
	adminRegistryToken := request.Header.Get("admin-registry-token")

	envAdminRegistryToken := os.Getenv("ADMIN_REGISTRY_TOKEN")

	switch {
	// a token is awaited
	case wasmArgs.AdminRegistryToken != "":
		if wasmArgs.AdminRegistryToken == adminRegistryToken {
			authorised = true
		} else {
			authorised = false
		}
	// a token is awaited
	case wasmArgs.AdminRegistryToken == "" && envAdminRegistryToken != "":
		if envAdminRegistryToken == adminRegistryToken {
			authorised = true
		} else {
			authorised = false
		}
	case wasmArgs.AdminRegistryToken == "" && envAdminRegistryToken == "":
		authorised = true
	default:
		authorised = false
	}

	return authorised
}

// CheckPrivateRegistryToken checks if the provided private registry token is authorized for the given request.
//
// Parameters:
// - request: the HTTP request object containing the private-registry-token header.
// - wasmArgs: the WasmArguments struct containing the private registry token.
//
// Return type:
// - bool: a boolean value indicating whether the token is authorized or not.
func CheckPrivateRegistryToken(request *http.Request, wasmArgs simplismTypes.WasmArguments) bool {
	var authorised bool = false
	// read the header private-registry-token
	privateRegistryToken := request.Header.Get("private-registry-token")

	envPrivateRegistryToken := os.Getenv("PRIVATE_REGISTRY_TOKEN")

	switch {
	// a token is awaited
	case wasmArgs.PrivateRegistryToken != "":
		if wasmArgs.PrivateRegistryToken == privateRegistryToken {
			authorised = true
		} else {
			authorised = false
		}
	// a token is awaited
	case wasmArgs.PrivateRegistryToken == "" && envPrivateRegistryToken != "":
		if envPrivateRegistryToken == privateRegistryToken {
			authorised = true
		} else {
			authorised = false
		}
	case wasmArgs.PrivateRegistryToken == "" && envPrivateRegistryToken == "":
		authorised = true
	default:
		authorised = false
	}

	return authorised
}