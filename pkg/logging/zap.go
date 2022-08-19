package logging

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	Level    zap.AtomicLevel
	DevMode  bool
	Encoding Encoding
}

type Encoding string

const (
	Console Encoding = "console"
	JSON             = "json"
)

func WithCfg(cfg *Config) (*zap.SugaredLogger, error) {

	builder := zap.NewProductionConfig()
	builder.Encoding = string(cfg.Encoding)
	builder.Level = cfg.Level
	builder.Development = cfg.DevMode

	defaultPath := "./logs/log.txt"

	_, err := os.Stat(defaultPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic(fmt.Sprintf("file %s does not exist", defaultPath))
		}
		panic(err.Error())
	}

	builder.OutputPaths = []string{defaultPath}

	logger, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
