package main

import (
	"datapoint/config"
	"datapoint/internal/app"
	"datapoint/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Must()
	zap.ReplaceGlobals(logger.Must(cfg.Logger))
	if err := app.Run(cfg); err != nil {
		zap.S().Fatal(err)
	}
}
