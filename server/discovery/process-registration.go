package discovery

import (
	"net/http"
	httpHelper "simplism/helpers/http"
	jsonHelper "simplism/helpers/json"
	data "simplism/server/data"
	simplismTypes "simplism/types"

	"go.etcd.io/bbolt"
)

func registerProcess(request *http.Request, db *bbolt.DB) (simplismTypes.SimplismProcess, error) {
	body := httpHelper.GetBody(request) // process information from simplism POST request
	// store the process information in the database
	simplismProcess, _ := jsonHelper.GetSimplismProcesseFromJSONBytes(body)
	err := data.SaveSimplismProcessToDB(db, simplismProcess)
	return simplismProcess, err
}
