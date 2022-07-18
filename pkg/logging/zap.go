package logging

import "go.uber.org/zap"

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

	logger, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
