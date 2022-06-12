package log

import "os"

func init() {
	if logger == nil {
		logger = NewDefaultLogger(os.Stdout)
	}
}

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func Debug(values ...interface{}) { logger.Debug(values...) }
func Info(values ...interface{})  { logger.Info(values...) }
func Warn(values ...interface{})  { logger.Warn(values...) }
func Error(values ...interface{}) { logger.Error(values...) }
func Fatal(values ...interface{}) { logger.Fatal(values...) } // In order to print the stack log

func Debugf(template string, values ...interface{}) { logger.Debugf(template, values...) }
func Infof(template string, values ...interface{})  { logger.Infof(template, values...) }
func Warnf(template string, values ...interface{})  { logger.Warnf(template, values...) }
func Errorf(template string, values ...interface{}) { logger.Errorf(template, values...) }
func Fatalf(template string, values ...interface{}) { logger.Fatalf(template, values...) } // In order to print the stack log
