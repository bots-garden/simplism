package server

// WasmArguments type
type WasmArguments struct {
	FilePath          string `yaml:"wasm-file,omitempty"`
	FunctionName      string `yaml:"wasm-function,omitempty"`
	HTTPPort          string `yaml:"http-port,omitempty"`
	Input             string `yaml:"input,omitempty"`
	LogLevel          string `yaml:"log-level,omitempty"`
	AllowHosts        string `yaml:"allow-hosts,omitempty"`
	AllowPaths        string `yaml:"allow-paths,omitempty"`
	EnvVars           string `yaml:"env,omitempty"`
	Config            string `yaml:"config,omitempty"`
	Wasi              bool   `yaml:"wasi,omitempty"`
	URL               string `yaml:"wasm-url,omitempty"`
	WasmURLAuthHeader string `yaml:"wasm-url-auth-header,omitempty"`
	//AuthHeaderName  string `yaml:"auth-header-name,omitempty"`
	//AuthHeaderValue string `yaml:"auth-header-value,omitempty"`
	CertFile         string `yaml:"cert-file,omitempty"`
	KeyFile          string `yaml:"key-file,omitempty"`
	AdminReloadToken string `yaml:"admin-reload-token,omitempty"`
	Discovery bool `yaml:"discovery,omitempty"`
	DiscoveryEndpoint string `yaml:"discovery-endpoint,omitempty"`
}

type SimplismProcess struct {
	PID          int
	FunctionName string
	FilePath     string
}
