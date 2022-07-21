package httptransport

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/binder"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"go.uber.org/zap"
	"net/http"
)

type runnerHandler struct {
	logger        *zap.SugaredLogger
	runnerService runner.Service
}

func NewRunnerHandler(logger *zap.SugaredLogger, runnerService runner.Service) runner.Transport {
	return &runnerHandler{logger: logger, runnerService: runnerService}
}

func (h *runnerHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/runner/", h.Register)
	r.DELETE("/api/runner/ban", h.Ban)
}

func (h *runnerHandler) Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var inp dto.RegisterRunnerDto
	err := binder.Bind(req.Body, &inp)
	if err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}

	err = h.runnerService.Register(inp)
	if err != nil {
		code, R := httpErrors.Response(err)
		responder.JSON(w, code, R)
		h.logger.Error(err.Error())
		return
	}

	w.WriteHeader(201)
	return
}

func (h *runnerHandler) Ban(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//TODO implement me
	panic("implement me")
}
