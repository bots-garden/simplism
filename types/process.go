package simplismTypes

import "time"

type SimplismProcess struct {
	PID          int       `json:"pid"`
	FunctionName string    `json:"functionName"`
	FilePath     string    `json:"filePath"`
	RecordTime   time.Time `json:"recordTime"`
	StartTime    time.Time `json:"startTime"`
}
