package httptransport

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/binder"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"go.uber.org/zap"
	"net/http"
)

type deliveryHandler struct {
	logger          *zap.SugaredLogger
	deliveryService delivery.Service
	telegramService telegram.Service
}

func NewDeliveryHandler(logger *zap.SugaredLogger, delivery delivery.Service, tg telegram.Service) delivery.Transport {
	return &deliveryHandler{logger: logger, deliveryService: delivery, telegramService: tg}
}

func (h *deliveryHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/delivery", h.CreateDelivery)
	r.POST("/api/delivery/status", h.Status)
	r.POST("/api/check", h.Check)
}

func (h *deliveryHandler) Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.logger.Info("call check")
	hdr := w.Header()

	hdr.Add("Content-Type", "octet/stream")
	hdr.Add("Connection", "keep-alive")

	var inp dto.CheckDto

	if err := binder.Bind(r.Body, &inp); err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}

	bytes, _ := json.Marshal(inp)

	str := string(bytes)
	fmt.Println(str)

}

func (h *deliveryHandler) CreateDelivery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	h.logger.Info("call create delivery")
	var payload dto.CreateDelivery

	if err := binder.Bind(req.Body, &payload); err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}
	createDto := dto.CreateDeliveryDatabaseDto{
		OrderID: payload.Order.OrderID,
		Pay:     payload.Order.Pay,
	}

	deliveryID, err := h.deliveryService.Create(createDto)
	if err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}
	h.logger.Debug("created delivery in database")
	//todo: mv template to templates, func to bot pkg
	telegramMsg := h.telegramService.FromTemplate(&payload)
	h.logger.Debug("formatted telegram template")

	//Data for telegram button callback query

	err = h.telegramService.Send(telegramMsg, deliveryID)
	if err != nil {

		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		//Delete delivery in database because telegram service could not send a message
		err = h.deliveryService.Delete(deliveryID)
		if err != nil {
			code, R = httpErrors.Response(err)
			responder.JSON(w, code, R)
			h.logger.Error(err.Error())
			return
		}
		return
	}
	h.logger.Debug("successfully sent telegram msg")
	w.WriteHeader(201)
	return
}

func (h *deliveryHandler) Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.logger.Info("call status")
	var inp dto.StatusOfDeliveryDto

	if err := binder.Bind(r.Body, &inp); err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}

	statuses, err := h.deliveryService.Status(inp)
	if err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}
	responder.JSON(w, http.StatusOK, responder.R{
		"result": statuses,
	})
	h.logger.Info("successfully sent statuses")
	return
}
