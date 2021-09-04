package utils

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"
)

type ILogger interface {
	Debug(ctx context.Context, pattern string, args ...interface{})
	Info(ctx context.Context, pattern string, args ...interface{})
	Warning(ctx context.Context, pattern string, args ...interface{})
	Error(ctx context.Context, err error)
}

func NewLogger(name, level string) (ILogger, error) {
	// 初始化內部 Logger
	internalLogger, err := newInternalLogger(name, level)
	if err != nil {
		return nil, xerrors.Errorf("啟動 Logger 失敗: %w", err)
	}

	return &logger{
		internalLogger: internalLogger,
	}, nil
}

type logger struct {
	internalLogger *zap.Logger
}

func (logger *logger) Debug(ctx context.Context, pattern string, args ...interface{}) {
	logger.internalLogger.Debug(fmt.Sprintf(pattern, args...))
}

func (logger *logger) Info(ctx context.Context, pattern string, args ...interface{}) {
	logger.internalLogger.Info(fmt.Sprintf(pattern, args...))
}

func (logger logger) Warning(ctx context.Context, pattern string, args ...interface{}) {
	logger.internalLogger.Warn(fmt.Sprintf(pattern, args...))
}

func (logger *logger) Error(ctx context.Context, err error) {
	logger.internalLogger.Error(fmt.Sprintf("%+v", err))
}

func newInternalLogger(name, level string) (*zap.Logger, error) {
	internalLoggerConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(getZapLevel(level)), // 日志级别
		Encoding: "console",                                // 輸出格式 console 或 json
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}, // 编码器配置
		InitialFields: map[string]interface{}{
			"service": name,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	internalLogger, err := internalLoggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, xerrors.Errorf("啟動內部 Logger 失敗: %w", err)
	}

	return internalLogger, nil
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case zapcore.DebugLevel.String(): // "debug"
		return zapcore.DebugLevel
	case zapcore.InfoLevel.String(): // "info"
		return zapcore.InfoLevel
	case zapcore.WarnLevel.String(): // "warn"
		return zapcore.WarnLevel
	case zapcore.ErrorLevel.String(): // "error"
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}
