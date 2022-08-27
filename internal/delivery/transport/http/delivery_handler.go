package httptransport

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sonyamoonglade/delivery-service/config"
	"github.com/sonyamoonglade/delivery-service/internal/delivery"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/binder"
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"github.com/sonyamoonglade/delivery-service/pkg/cli"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	"github.com/sonyamoonglade/delivery-service/pkg/formatter"
	"github.com/sonyamoonglade/notification-service/pkg/httpRes"
	"go.uber.org/zap"
)

type deliveryHandler struct {
	logger          *zap.SugaredLogger
	deliveryService delivery.Service
	extractFmt      formatter.ExtractFormatter
	bot             bot.Bot
}

func NewDeliveryHandler(logger *zap.SugaredLogger, delivery delivery.Service, extractFormatter formatter.ExtractFormatter, bot bot.Bot) delivery.Transport {
	return &deliveryHandler{logger: logger, deliveryService: delivery, extractFmt: extractFormatter, bot: bot}
}

func (h *deliveryHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/delivery", h.CreateDelivery)
	r.POST("/api/delivery/status", h.Status)
	r.POST("/api/check", h.Check)

}

func (h *deliveryHandler) Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	h.logger.Debug("call check")
	hdr := w.Header()

	hdr.Add("Connection", "keep-alive")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	doneErrCh := make(chan error)

	var inp dto.CheckDto
	if err := binder.Bind(r.Body, &inp); err != nil {
		httpErrors.ResponseAndLog(h.logger, w, err)
		return
	}

	dtoForCli := dto.CheckDtoForCli{
		Data: inp,
	}

	//Imitate timeout
	go func() {
		time.Sleep(check.CheckWriteTimeout)
		cancel()
	}()

	//Invoke write in goroutine
	go func() {
		err := h.deliveryService.WriteCheck(dtoForCli)
		if err != nil {
			doneErrCh <- err
			return
		}
		doneErrCh <- nil
	}()
	//Block initial routine with select case
	select {
	case <-ctx.Done():
		h.logger.Errorf("Failed with timeout. %s", ctx.Err())
		httpErrors.ResponseAndLog(h.logger, w, cli.TimeoutError)
		return
	case err := <-doneErrCh:
		if err != nil {
			httpErrors.ResponseAndLog(h.logger, w, err)
			return
		}
		//If previous operations were ok, set header
		hdr.Add("Content-Type", "octet/stream")

		//Copy check file bytes to ResponseWriter
		err = h.deliveryService.CopyFromCheck(w)
		if err != nil {
			httpErrors.ResponseAndLog(h.logger, w, err)
			return
		}

		h.logger.Debug("copy file to response writer success")
		return
	}

}

func (h *deliveryHandler) CreateDelivery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	h.logger.Debug("call create delivery")
	var payload dto.CreateDelivery

	if err := binder.Bind(req.Body, &payload); err != nil {
		httpErrors.ResponseAndLog(h.logger, w, err)
		return
	}
	createDto := dto.CreateDeliveryDatabaseDto{
		OrderID: payload.Order.OrderID,
		Pay:     payload.Order.Pay,
	}

	deliveryID, err := h.deliveryService.Create(createDto)
	if err != nil {
		httpErrors.ResponseAndLog(h.logger, w, err)
		return
	}
	h.logger.Debug("created delivery in database")

	telegramMsg := h.extractFmt.FormatTemplate(&payload, config.TempOffset)
	h.logger.Debug("formatted telegram template")

	//Bot produced an error while sending a message
	err = h.bot.PostDeliveryMessage(telegramMsg, deliveryID)
	if err != nil {
		//Delete delivery in database because telegram service could not send a message
		err = h.deliveryService.Delete(deliveryID)
		if err != nil {
			httpErrors.ResponseAndLog(h.logger, w, err)
			return
		}

		httpErrors.ResponseAndLog(h.logger, w, err)
		return
	}
	h.logger.Debug("successfully sent telegram msg")
	httpRes.Created(w)
}

func (h *deliveryHandler) Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.logger.Info("request for statuses")
	var inp dto.StatusOfDeliveryDto

	if err := binder.Bind(r.Body, &inp); err != nil {
		httpErrors.ResponseAndLog(h.logger, w, err)
		return
	}

	statuses, err := h.deliveryService.Status(inp)
	if err != nil {
		httpErrors.ResponseAndLog(h.logger, w, err)
		return
	}
	httpRes.Json(h.logger, w, http.StatusOK, httpRes.JSON{
		"result": statuses,
	})
	h.logger.Debug("sent")
	return
}
