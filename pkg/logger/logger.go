package logger

import (
	"datapoint/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Must(cfg config.Logger) *zap.Logger {
	c := zap.NewDevelopmentConfig()

	if cfg.Production {
		c = zap.NewProductionConfig()
	}

	c.DisableStacktrace = cfg.DisableStacktrace

	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zap.Must(c.Build())
}
