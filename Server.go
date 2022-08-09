package tgdelivery

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sonyamoonglade/delivery-service/config"
	"net/http"
	"time"
)

func NewServerWithConfig(cfg *config.App, h *httprouter.Router) *http.Server {
	s := http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        h,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 10 << 15,
	}

	return &s
}
