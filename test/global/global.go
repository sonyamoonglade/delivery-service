package global

import (
	"go.uber.org/zap"
)

func InitLogger() *zap.SugaredLogger {
	prod, _ := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "json",
	}.Build()
	logger := prod.Sugar()
	return logger
}
