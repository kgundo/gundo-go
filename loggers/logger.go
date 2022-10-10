package loggers

import (
	"github.com/sirupsen/logrus"
)

type Event struct {
	id      int
	message string
}

type StandardLogger struct {
	*logrus.Logger
}

func Logger() *StandardLogger {
	var baseLogger = logrus.New()
	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}
	return standardLogger
}

var (
	commonExceptionMessage = Event{1, "Exception: %s"}
	invalidArgMessage      = Event{10, "Invalid arg: %s"}
	invalidArgValueMessage = Event{20, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{30, "Missing arg: %s"}
	commonInfoMessage      = Event{1, "Info: %s"}
)

func (l *StandardLogger) CommonException(err string) {
	l.Errorf(commonExceptionMessage.message, err)
}

func (l *StandardLogger) InvalidArg(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

func (l *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

func (l *StandardLogger) MissingArg(argumentName string) {
	l.Errorf(missingArgMessage.message, argumentName)
}

func (l *StandardLogger) CommonInfo(argumentName string) {
	l.Infof(commonInfoMessage.message, argumentName)
}
