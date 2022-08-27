package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Transport interface {
	RegisterRoutes(r *httprouter.Router)
	Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params)
	Ban(w http.ResponseWriter, req *http.Request, _ httprouter.Params)
}
