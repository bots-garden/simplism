package spawn

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func subHandlerFallAsleepByPid(request *http.Request, response http.ResponseWriter) {
	spid := chi.URLParam(request, "pid")
	pid, err := strconv.Atoi(spid)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("pid not present"))
	} else {
		// kill the process
		foundProcess, errAsleep := fallAsleepProcess(pid)
		if errAsleep != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(errAsleep.Error()))
		}
		response.WriteHeader(http.StatusOK)
		response.Write([]byte(foundProcess.ServiceName + "[" + strconv.Itoa(foundProcess.PID) + "]" + " asleep"))
	}
}
