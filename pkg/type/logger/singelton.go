package logger

import (
	"go.uber.org/zap"

	"architecture_go/pkg/type/context"
)

var log *Logger

func init() {
	if log == nil {
		newLogger, err := new()
		if err != nil {
			panic(err)
		}

		log = newLogger
	}
}

func GetLogger() *zap.Logger {
	return log.logger
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func DebugWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	log.DebugWithContext(ctx, msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func InfoWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	log.InfoWithContext(ctx, msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func WarnWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	log.WarnWithContext(ctx, msg, fields...)
}

func Error(msg interface{}, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func ErrorWithContext(ctx context.Context, err error, fields ...zap.Field) error {
	return log.ErrorWithContext(ctx, err, fields...)
}

func Fatal(msg interface{}, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func FatalWithContext(ctx context.Context, err error, fields ...zap.Field) error {
	return log.FatalWithContext(ctx, err, fields...)
}
