package internal

import (
	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error initializing zap logger: " + err.Error())
	}
	return logger
}
