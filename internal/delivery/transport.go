package delivery

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Transport interface {
	RegisterRoutes(r *httprouter.Router)
	CreateDelivery(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}
