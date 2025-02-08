package logger

import (
	"github.com/DhruvikDonga/mock-ses/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) *zap.SugaredLogger {
	var logLevel int
	if cfg.Build == "release" {
		logLevel = int(zapcore.InfoLevel)
	} else {
		logLevel = int(zapcore.DebugLevel)
	}
	zapCfg := zap.Config{
		Encoding:         "console", // values can be json or console
		Level:            zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		OutputPaths:      []string{"stdout"}, // or specify a file path
		ErrorOutputPaths: []string{"stdout"}, // or specify a file path for error logs
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "Message",
			LevelKey:         "Loglevel",
			TimeKey:          "Timestamp",
			NameKey:          "Logger",
			CallerKey:        "Files",
			StacktraceKey:    "Stacktrace",
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeCaller:     nil,
			ConsoleSeparator: " ",
		},
	}
	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
