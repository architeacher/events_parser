package logger

import (
	baseLogger "log"
	"os"
)

type Logger struct {
	baseLogger *baseLogger.Logger
	formatter  *Formatter
}

type Context interface{}

type Level uint

const (
	EMERGENCY Level = 8 - iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
	TRACE
)

func NewLogger() *Logger {

	return &Logger{
		baseLogger: baseLogger.New(os.Stderr, "", baseLogger.LstdFlags),
		formatter:  NewFormatter(),
	}
}

func (self *Logger) Emergency(context ...Context) {
	self.baseLogger.Fatal(self.formatter.format(context))
}

func (self *Logger) Alert(context ...Context) {
	self.baseLogger.Fatal(self.formatter.format(context))
}

func (self *Logger) Critical(context ...Context) {
	self.baseLogger.Fatal(self.formatter.format(context))
}

func (self *Logger) Error(context ...Context) {
	self.baseLogger.Println(self.formatter.format(context))
}

func (self *Logger) Warning(context ...Context) {
	self.baseLogger.Println(self.formatter.format(context))
}

func (self *Logger) Notice(context ...Context) {
	self.baseLogger.Println(self.formatter.format(context))
}

func (self *Logger) Info(context ...Context) {
	self.baseLogger.Println(self.formatter.format(context))
}

func (self *Logger) Debug(context ...Context) {
	self.baseLogger.Println(self.formatter.format(context))
}

func (self *Logger) Trace(context ...Context) {
	self.baseLogger.Println(self.formatter.format(context))
}

func levelAsString(l Level) string {
	switch l {
	case EMERGENCY:
		return "EMERGENCY"
	case ALERT:
		return "ALERT"
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	}

	return "TRACE"
}
