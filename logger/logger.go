package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger zap.SugaredLogger

func init() {
	conf := zap.NewProductionConfig()
	conf.Encoding = "console"
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	unsugared, err := conf.Build()
	if err != nil {
		log.Fatalln(err)
	}
	Logger = *unsugared.Sugar()

}
