package app

import "github.com/sirupsen/logrus"

type Logger interface {
	SetField(name, value string)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type logger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

func NewLogger(log *logrus.Logger, fields logrus.Fields) Logger {
	return &logger{
		logger: log,
		fields: fields,
	}
}

func (log *logger) SetField(name, value string) {
	log.fields[name] = value
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.tagged().Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.tagged().Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.tagged().Warnf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.tagged().Errorf(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.tagged().Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.tagged().Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.tagged().Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.tagged().Error(args...)
}

func (l *logger) tagged() *logrus.Entry {
	return l.logger.WithFields(l.fields)
}
