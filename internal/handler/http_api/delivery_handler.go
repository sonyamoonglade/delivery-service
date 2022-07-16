package apihandler

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	api_erros "github.com/sonyamoonglade/delivery-service/internal/errors/api"
	tg_errors "github.com/sonyamoonglade/delivery-service/internal/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/internal/handler/dto"
	"github.com/sonyamoonglade/delivery-service/internal/service"
	"github.com/sonyamoonglade/delivery-service/pkg/binder"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"go.uber.org/zap"
	"net/http"
)

type DeliveryHandler struct {
	logger          *zap.Logger
	deliveryService service.Delivery
	telegramService service.Telegram
}

func NewDeliveryHandler(logger *zap.Logger, service service.Delivery, tgservice service.Telegram) *DeliveryHandler {
	return &DeliveryHandler{logger: logger, deliveryService: service, telegramService: tgservice}
}

func (h *DeliveryHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/delivery", h.CreateDelivery)

}

func (h *DeliveryHandler) CreateDelivery(r http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	payload, err := binder.Bind(req.Body)

	if err != nil {
		var e *binder.BindingError
		if errors.As(err, &e) {
			responder.JSON(r, responder.R{
				"message": err.Error(),
			}, http.StatusBadRequest)
			h.logger.Error(err.Error())
			return
		}
		responder.JSON(r, responder.R{
			"message": api_erros.InternalServerError,
		}, http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	createDto := &dto.CreateDeliveryDto{
		OrderID: payload.Order.OrderID,
		Pay:     payload.Order.Pay,
	}

	err = h.deliveryService.Create(createDto)
	if err != nil {

		return
	}

	err = h.telegramService.Send(payload)
	if err != nil {

		var e tg_errors.TelegramError
		if errors.As(err, &e) {
			responder.JSON(r, responder.R{
				"message": err.Error(),
			}, http.StatusInternalServerError)
			h.logger.Error(err.Error())
			return
		}

		responder.JSON(r, responder.R{
			"message": api_erros.InternalServerError,
		}, http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	return
}
