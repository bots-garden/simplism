package server

import "time"

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
	Discovery           bool   `yaml:"discovery,omitempty"`
	DiscoveryEndpoint   string `yaml:"discovery-endpoint,omitempty"`
	AdminDiscoveryToken string `yaml:"admin-discovery-token,omitempty"`
}

type SimplismProcess struct {
	PID          int       `json:"pid"`
	FunctionName string    `json:"functionName"`
	FilePath     string    `json:"filePath"`
	RecordTime   time.Time `json:"recordTime"`
	StartTime    time.Time `json:"startTime"`
}
