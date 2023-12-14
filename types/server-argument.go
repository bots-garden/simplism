package simplismTypes

// WasmArguments type
type WasmArguments struct {
	FilePath            string `yaml:"wasm-file,omitempty"`
	FunctionName        string `yaml:"wasm-function,omitempty"`
	HTTPPort            string `yaml:"http-port,omitempty"`
	Input               string `yaml:"input,omitempty"`
	LogLevel            string `yaml:"log-level,omitempty"`
	AllowHosts          string `yaml:"allow-hosts,omitempty"`
	AllowPaths          string `yaml:"allow-paths,omitempty"`
	EnvVars             string `yaml:"env,omitempty"`
	Config              string `yaml:"config,omitempty"`
	Wasi                bool   `yaml:"wasi,omitempty"`
	URL                 string `yaml:"wasm-url,omitempty"`
	WasmURLAuthHeader   string `yaml:"wasm-url-auth-header,omitempty"`
	CertFile            string `yaml:"cert-file,omitempty"`
	KeyFile             string `yaml:"key-file,omitempty"`
	AdminReloadToken    string `yaml:"admin-reload-token,omitempty"`
	ServiceDiscovery    bool   `yaml:"service-discovery,omitempty"`
	DiscoveryEndpoint   string `yaml:"discovery-endpoint,omitempty"`
	AdminDiscoveryToken string `yaml:"admin-discovery-token,omitempty"`
	SpawnMode           bool   `yaml:"spawn-mode,omitempty"`
	AdminSpawnToken     string `yaml:"admin-spawn-token,omitempty"`
	Information         string `yaml:"information,omitempty"`
	ServiceName         string `yaml:"service-name,omitempty"`
	StoreMode           bool   `yaml:"store-mode,omitempty"`
	AdminStoreToken     string `yaml:"admin-store-token,omitempty"`
}

/*
	When adding a field to the structure, you need to update these files:
	- cmds/listen.go
    - server/handler-spawn.go
*/