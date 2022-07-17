package httptransport

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sonyamoonglade/delivery-service/internal/runner"
	"go.uber.org/zap"
	"net/http"
)

type runnerHandler struct {
	logger        *zap.Logger
	runnerService runner.Service
}

func NewRunnerHandler(logger *zap.Logger, runnerService runner.Service) runner.Transport {
	return &runnerHandler{logger: logger, runnerService: runnerService}
}

func (h *runnerHandler) RegisterRoutes(r *httprouter.Router) {

	r.POST("/api/runner/", h.Register)
	r.GET("/api/runner/isRunner", h.IsRunner)
	r.DELETE("/api/runner/ban", h.Ban)
}

func (h *runnerHandler) IsRunner(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//TODO implement me
	panic("implement me")
}

func (h *runnerHandler) Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	//out, err := binder.Bind(req.Body, dto.RegisterRunnerDto{})
	//if err != nil {
	//
	//	code, R := httpErrors.Response(err)
	//	responder.JSON(w, code, R)
	//	h.logger.Error(err.Error())
	//	return
	//}
	//w.Write([]byte("hey"))
	//return
}

func (h *runnerHandler) Ban(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//TODO implement me
	panic("implement me")
}
