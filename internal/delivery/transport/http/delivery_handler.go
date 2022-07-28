package httptransport

import (
	"context"
	"github.com/julienschmidt/httprouter"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/internal/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/binder"
	"github.com/sonyamoonglade/delivery-service/pkg/cli"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"go.uber.org/zap"
	"net/http"
	"time"
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

	hdr.Add("Connection", "keep-alive")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	doneCh := make(chan error)

	var inp dto.CheckDto
	if err := binder.Bind(r.Body, &inp); err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}

	dtoForCli := dto.CheckDtoForCli{
		Data: inp,
	}

	//Imitate timeout
	go func() {
		time.Sleep(tgdelivery.CheckWriteTimeout)
		cancel()
	}()

	go func() {
		err := h.deliveryService.WriteCheck(dtoForCli)
		if err != nil {
			doneCh <- err
			return
		}
		doneCh <- nil
	}()

	select {
	case <-ctx.Done():
		h.logger.Errorf("Failed with timeout. %s", ctx.Err())
		code, R := httpErrors.Response(cli.TimeoutError)
		responder.JSON(w, code, R)
		return
	case err := <-doneCh:
		if err != nil {
			code, R := httpErrors.Response(err)
			h.logger.Error(err.Error())
			responder.JSON(w, code, R)
			return
		}
		//If previous operations were ok, set header
		hdr.Add("Content-Type", "octet/stream")

		//Copy check file bytes to ResponseWriter
		err = h.deliveryService.ReadFromCheck(w)
		if err != nil {
			code, R := httpErrors.Response(err)
			h.logger.Error(err.Error())
			responder.JSON(w, code, R)
			return
		}

		h.logger.Info("copy file to response writer success")
		return
	}

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
	telegramMsg := h.telegramService.FormatTemplate(&payload)
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
