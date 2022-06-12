package log

import (
	"fmt"
	"io"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

func NewDefaultLogger(output io.Writer) Logger {
	return &DefaultLogger{writer: output}
}

type DefaultLogger struct {
	name   string
	writer io.Writer
}

func (l *DefaultLogger) Log(level string, msg string) {
	l.writer.Write([]byte("["))
	l.writer.Write([]byte(level))
	l.writer.Write([]byte("] "))
	l.writer.Write([]byte(msg + "\n"))
}

func (l *DefaultLogger) Debug(values ...interface{}) { l.Log("DEG", fmt.Sprint(values...)) }
func (l *DefaultLogger) Info(values ...interface{})  { l.Log("INF", fmt.Sprint(values...)) }
func (l *DefaultLogger) Warn(values ...interface{})  { l.Log("WAN", fmt.Sprint(values...)) }
func (l *DefaultLogger) Error(values ...interface{}) { l.Log("ERR", fmt.Sprint(values...)) }
func (l *DefaultLogger) Fatal(values ...interface{}) { l.Log("FAT", fmt.Sprint(values...)) }

func (l *DefaultLogger) Debugf(template string, values ...interface{}) {
	l.Log("DEG", fmt.Sprintf(template, values...))
}
func (l *DefaultLogger) Infof(template string, values ...interface{}) {
	l.Log("INF", fmt.Sprintf(template, values...))
}
func (l *DefaultLogger) Warnf(template string, values ...interface{}) {
	l.Log("WAN", fmt.Sprintf(template, values...))
}
func (l *DefaultLogger) Errorf(template string, values ...interface{}) {
	l.Log("ERR", fmt.Sprintf(template, values...))
}
func (l *DefaultLogger) Fatalf(template string, values ...interface{}) {
	l.Log("FAT", fmt.Sprintf(template, values...))
}
