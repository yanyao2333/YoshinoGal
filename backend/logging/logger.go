package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var log *zap.SugaredLogger

//var log = &logrus.Logger{
//	Out:       os.Stdout,
//	Formatter: new(logrus.TextFormatter),
//	Hooks:     make(logrus.LevelHooks),
//	Level:     logrus.DebugLevel,
//}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.WriteSyncer(os.Stdout)
}

func initLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	log = logger.Sugar()
}

// GetLogger 获取日志对象
func GetLogger() *zap.SugaredLogger {
	if log == nil {
		initLogger()
	}
	return log
}
