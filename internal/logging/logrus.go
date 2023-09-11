package logging

import (
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger(level logrus.Level) *LogrusLogger {
	logger := logrus.New()
	logger.Level = level
	return &LogrusLogger{
		logger: logger,
	}
}

func (l LogrusLogger) Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Debug(msg)
	} else {
		l.logger.Debug(msg, args)
	}

}

func (l LogrusLogger) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(msg)
	} else {
		l.logger.Info(msg, args)
	}

}

func (l LogrusLogger) Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Warn(msg)
	} else {
		l.logger.Warn(msg, args)
	}
}

func (l LogrusLogger) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Error(msg)
	} else {
		l.logger.Error(msg, args)
	}
}

func (l LogrusLogger) Fatal(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Fatal(msg)
	} else {
		l.logger.Fatal(msg, args)
	}
}
