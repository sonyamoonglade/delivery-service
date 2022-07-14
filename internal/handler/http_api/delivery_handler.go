package apihandler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/service"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type DeliveryHandler struct {
	logger    *zap.Logger
	service   service.Delivery
	tgservice service.Telegram
}

func NewDeliveryHandler(logger *zap.Logger, service service.Delivery, tgservice service.Telegram) *DeliveryHandler {
	return &DeliveryHandler{logger: logger, service: service, tgservice: tgservice}
}

func (h *DeliveryHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/delivery", h.CreateDelivery)

}

func (h *DeliveryHandler) CreateDelivery(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var p tgdelivery.Payload

	b, _ := io.ReadAll(req.Body)
	json.Unmarshal(b, &p)

	rw.Write([]byte("Hello world"))
	h.tgservice.Send(&p)

	return
}
