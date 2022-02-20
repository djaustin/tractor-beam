package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger zap.SugaredLogger
var logLevel zap.AtomicLevel

func SetLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}

func init() {
	conf := zap.NewProductionConfig()

	logLevel = zap.NewAtomicLevel()

	conf.Level = logLevel
	conf.Encoding = "console"
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	unsugared, err := conf.Build()
	if err != nil {
		log.Fatalln(err)
	}
	Logger = *unsugared.Sugar()

}
