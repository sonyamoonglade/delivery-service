package runner

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Transport interface {
	RegisterRoutes(r *httprouter.Router)
	Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params)
	Ban(w http.ResponseWriter, req *http.Request, _ httprouter.Params)
}
