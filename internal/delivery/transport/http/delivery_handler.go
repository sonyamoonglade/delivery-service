package http

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/binder"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/api"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"go.uber.org/zap"
	"net/http"
)

type DeliveryHandler struct {
	logger          *zap.Logger
	deliveryService delivery.Delivery
	telegramService telegram.Telegram
}

func NewDeliveryHandler(logger *zap.Logger, delivery delivery.Delivery, tg telegram.Telegram) delivery.Transport {
	return &DeliveryHandler{logger: logger, deliveryService: delivery, telegramService: tg}
}

func (h *DeliveryHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/delivery", h.CreateDelivery)

}

func (h *DeliveryHandler) CreateDelivery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	payload, err := binder.Bind(req.Body)
	if err != nil {
		var e *binder.BindingError
		if errors.As(err, &e) {
			responder.JSON(w, responder.R{
				"message": err.Error(),
			}, http.StatusBadRequest)
			h.logger.Error(err.Error())
			return
		}
		responder.JSON(w, responder.R{
			"message": api_erros.InternalServerError,
		}, http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	createDto := &dto.CreateDeliveryDto{
		OrderID: payload.Order.OrderID,
		Pay:     payload.Order.Pay,
	}

	deliveryID, err := h.deliveryService.Create(createDto)
	if err != nil {

		return
	}
	h.logger.Debug("created delivery in database")

	telegramMsg := h.telegramService.FromTemplate(payload)
	h.logger.Debug("formatted telegram template")

	err = h.telegramService.Send(telegramMsg)
	if err != nil {

		var e tg_errors.TelegramError
		if errors.As(err, &e) {
			responder.JSON(w, responder.R{
				"message": err.Error(),
			}, http.StatusInternalServerError)
			h.logger.Error(err.Error())
			return
		}

		responder.JSON(w, responder.R{
			"message": api_erros.InternalServerError,
		}, http.StatusInternalServerError)
		h.logger.Error(err.Error())

		err = h.deliveryService.Delete(deliveryID)
		if err != nil {

			var e api_erros.InvalidDelivery
			if errors.As(err, &e) {
				responder.JSON(w, responder.R{
					"message": err.Error(),
				}, http.StatusBadRequest)
				h.logger.Debug(err.Error())
				return
			}

			responder.JSON(w, responder.R{
				"message": api_erros.InternalServerError,
			}, http.StatusInternalServerError)
			h.logger.Error(err.Error())

			return
		}

		return
	}
	h.logger.Debug("successfully sent telegram msg")

	return
}
