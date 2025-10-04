package logs

import (
	"go.uber.org/zap"
)

func NewZapLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("Can't initialize zap logger: " + err.Error())
		panic(err)
	}

	logger.Info("Zap logger success")

	return logger.Sugar()
}
