package processes

import (
	simplismTypes "simplism/types"
	"time"
)

var currentSimplismProcess = simplismTypes.SimplismProcess{}

func GetCurrentSimplismProcess() simplismTypes.SimplismProcess {
	return currentSimplismProcess
}

func SetCurrentProcessPID(pid int) {
	currentSimplismProcess.PID = pid
}

func SetCurrentProcessFilePath(filePath string) {
	currentSimplismProcess.FilePath = filePath
}

func SetCurrentProcessFunctionName(functionName string) {
	currentSimplismProcess.FunctionName = functionName
}

func SetCurrentProcessHTTPPort(httpPort string) {
	currentSimplismProcess.HTTPPort = httpPort
}

func SetCurrentProcessInformation(information string) {
	currentSimplismProcess.Information = information
}

func SetCurrentProcessServiceName(serviceName string) {
	currentSimplismProcess.ServiceName = serviceName
}

func SetCurrentProcessStartTime(startTime time.Time) {
	currentSimplismProcess.StartTime = startTime
}

func SetCurrentProcessData(data simplismTypes.SimplismProcess) {
    currentSimplismProcess = data
}
