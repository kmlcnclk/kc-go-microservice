package log

import (
	"fmt"
	"hash/fnv"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func hashToColorCode(name string) int {
	h := fnv.New32a()
	h.Write([]byte(name))
	return int(33 + (h.Sum32() % 198))
}

func coloredNameEncoder(name string, enc zapcore.PrimitiveArrayEncoder) {
	colorCode := hashToColorCode(name)
	colored := fmt.Sprintf("\033[1;38;5;%dm%s\033[0m", colorCode, name)
	enc.AppendString(colored)
}

func Init(serviceName string) {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     coloredNameEncoder,
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	logger = zap.Must(config.Build()).Named(serviceName)
	zap.ReplaceGlobals(logger)
}
