package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	logLevel  = "LOG_LEVEL"
	logOutput = "LOG_OUTPUT"
)

var (
	log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput()},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
			CallerKey:    "caller",
		},
	}
	var err error
	if log, err = logConfig.Build(zap.AddCallerSkip(1)); err != nil {
		panic(err)
	}
}

func getOutput() string {
	output := strings.TrimSpace(os.Getenv(logOutput))
	if output == "" {
		return "stdout"
	}

	return output
}

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(logLevel))) {
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.InfoLevel
	}

}

func Info(msg string, tags ...zap.Field) {
	log.Info(msg, tags...)
	log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.Error(err))
	log.Error(msg, tags...)
	log.Sync()
}

func Fatal(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.Error(err))
	log.Fatal(msg, tags...)
	log.Sync()
}
