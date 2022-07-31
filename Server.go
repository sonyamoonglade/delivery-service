package tgdelivery

import (
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func NewServerWithConfig(cfg *viper.Viper, h *httprouter.Router) *http.Server {
	s := http.Server{
		Addr:           ":" + cfg.GetString("app.port"),
		Handler:        h,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 10 << 15,
	}

	return &s
}
