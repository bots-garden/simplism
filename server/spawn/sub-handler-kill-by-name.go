package spawn

import (
	"net/http"
	"simplism/server/discovery"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func subHandlerKillByName(request *http.Request, response http.ResponseWriter) {

	serviceName := chi.URLParam(request, "name")

	foundProcess, err := discovery.NotifyProcesseInformation(serviceName)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("service not found"))
	}
	// kill the process
	_, errKill := killProcess(foundProcess.PID)
	if errKill != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.WriteHeader(http.StatusOK)
	response.Write([]byte(foundProcess.ServiceName + "[" + strconv.Itoa(foundProcess.PID) + "]" + " killed"))
}
