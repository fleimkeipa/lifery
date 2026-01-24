package logger

import (
	"os"

	"github.com/fleimkeipa/lifery/model"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global sugared logger
var Log *zap.SugaredLogger

// Init initializes the global logger based on the environment
func Init() {
	stage := os.Getenv("STAGE")
	var config zap.Config

	opts := []zap.Option{}

	switch stage {
	case model.StageProd:
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.TimeKey = "timestamp"
	default:
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		opts = append(opts, zap.AddCallerSkip(1))
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}

	// Build the logger
	l, err := config.Build(opts...)
	if err != nil {
		panic(err)
	}

	Log = l.Sugar()
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
