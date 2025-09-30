package logs

import (
	"os"

	"go.uber.org/zap"
)

func GetLogger() zap.Logger {
	_, err := os.ReadDir("logs")
	if err != nil {
		os.Mkdir("logs", 0755)
	}
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "logs/app.log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		InitialFields:    map[string]interface{}{},
	}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	return *logger
}
