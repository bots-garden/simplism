package simplismTypes

import "time"

type SimplismProcess struct {
	PID          int       `json:"pid"`
	FunctionName string    `json:"functionName"`
	FilePath     string    `json:"filePath"`
	RecordTime   time.Time `json:"recordTime"`
	StartTime    time.Time `json:"startTime"`
	HTTPPort     string    `json:"httpPort"`
	Information  string    `json:"information"` // not used, but soon
	ServiceName  string    `json:"serviceName"`
	//Host         string    `json:"host"` // how to set this?
}
