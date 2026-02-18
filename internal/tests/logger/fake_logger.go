package logger

import application "health-checker/internal/application/logger"

type FakeLogger struct{}

func NewFakeLogger() *FakeLogger {
	return &FakeLogger{}
}

func (l *FakeLogger) Info(msg string, fields ...application.Field)  {}
func (l *FakeLogger) Warn(msg string, fields ...application.Field)  {}
func (l *FakeLogger) Error(msg string, fields ...application.Field) {}
func (l *FakeLogger) Debug(msg string, fields ...application.Field) {}
func (l *FakeLogger) Fatal(msg string, fields ...application.Field) {}
