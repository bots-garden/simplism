package spawn

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func subHandlerKillByPid(request *http.Request, response http.ResponseWriter) {
	spid := chi.URLParam(request, "pid")
	pid, err := strconv.Atoi(spid)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("pid not present"))
	} else {
		// kill the process
		foundProcess, errKill := killProcess(pid)
		if errKill != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(errKill.Error()))
		}
		response.WriteHeader(http.StatusOK)
		response.Write([]byte(foundProcess.ServiceName + "[" + strconv.Itoa(foundProcess.PID) + "]" + " killed"))
	}
}
