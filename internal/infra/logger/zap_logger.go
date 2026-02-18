package logger

import (
	application "health-checker/internal/application/logger"

	"go.uber.org/zap"
)

type ZapLogger struct {
	log *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &ZapLogger{log: logger}, nil
}

func (z *ZapLogger) Info(msg string, fields ...application.Field) {
	z.log.Info(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...application.Field) {
	z.log.Warn(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...application.Field) {
	z.log.Error(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Debug(msg string, fields ...application.Field) {
	z.log.Debug(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...application.Field) {
	z.log.Fatal(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) toZapFields(fields []application.Field) []zap.Field {
	var zapFields []zap.Field
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
