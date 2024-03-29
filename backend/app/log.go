package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//var log = &logrus.Logger{
//	Out:       os.Stdout,
//	Formatter: new(logrus.TextFormatter),
//	Hooks:     make(logrus.LevelHooks),
//	Level:     logrus.DebugLevel,
//}

var log *zap.SugaredLogger

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.WriteSyncer(os.Stdout)
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	log = logger.Sugar()
}
