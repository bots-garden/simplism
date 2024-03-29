package simplismTypes

import "time"

type SimplismProcess struct {
	PID          int       `json:"pid"`
	FunctionName string    `json:"functionName"`
	FilePath     string    `json:"filePath"`
	RecordTime   time.Time `json:"recordTime"`
	StartTime    time.Time `json:"startTime"`
	StopTime     time.Time `json:"stopTime"`
	HTTPPort     string    `json:"httpPort"`
	Information  string    `json:"information"` // not used, but soon
	ServiceName  string    `json:"serviceName"`
	Host         string    `json:"host"` // how to set this?
	Asleep       bool      `json:"asleep"`
}

//TODO add default value to Host
//? when this is set?
