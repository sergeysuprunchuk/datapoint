package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Must() *zap.Logger {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.Must(c.Build())
}
