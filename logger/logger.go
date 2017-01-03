package logger

import (
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

//Logger is the wrapper or logrus logger.
// var Logger = logrus.WithFields(logrus.Fields{
// 	"application": "bookmark",
// })

var Logger *logrus.Entry

func init() {
	instance := logrus.New()
	instance.Formatter = new(logrus.JSONFormatter)
	instance.Formatter = new(logrus.TextFormatter)
	instance.Level = logrus.InfoLevel
	instance.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  "./log/info.log",
		logrus.ErrorLevel: "./log/error.log",
		logrus.PanicLevel: "./log/panic.log",
		logrus.FatalLevel: "./log/fatal.log",
	}))

	Logger = instance.WithFields(logrus.Fields{
		"application": "bookmark",
	})
	Logger = decorateRuntimeContext(Logger)
}

func decorateRuntimeContext(logger *logrus.Entry) *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	}
	return logger

}
