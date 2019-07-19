package slog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	MetricLog = "csp/sentinel-metirc.log"
	BlockLog  = "csp/sentinel-block.log"
	Record    = "csp/sentinel-record.log"
)

var logs map[string]*zap.Logger

func init() {
	logs = make(map[string]*zap.Logger)
	addLog(MetricLog, zapcore.InfoLevel)
	addLog(BlockLog, zapcore.InfoLevel)
	addLog(Record, zapcore.InfoLevel)
}

func addLog(filePath string, level zapcore.Level) {
	jLoger := &lumberjack.Logger{
		Filename: filePath,
		MaxSize:  100, //MB
	}
	w := zapcore.AddSync(jLoger)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		w,
		level,
	)

	log := zap.New(core)
	log = log.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	logs[filePath] = log
}

func GetLog(filePath string) *zap.Logger {
	return logs[filePath]
}
