package router

import 	"github.com/go-chi/chi/v5"


var router = chi.NewRouter()

func GetRouter() *chi.Mux {
	return router
}
