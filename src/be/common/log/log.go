package log

import (
	"be/options"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Fields = logrus.Fields

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func InitLog() {
	// 设置日志输出格式
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	// 设置默认的日志级别
	switch strings.ToUpper(options.Options.LogLevel) {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}

	// 设置输出路径
	f, err := os.OpenFile(options.Options.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("无法打开日志文件")
	}
	logrus.SetOutput(f)
}

func Fatal(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Fatal(args...)
}

func Debugln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Debugln(args...)
}

func Infoln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Infoln(args...)
}

func Warnln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Warnln(args...)
}

func Errorln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Errorln(args...)
}

func Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("gofile", fmt.Sprintf("[%s:%d] ", file, line)).Errorf(format, args...)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	fields["gofile"] = fmt.Sprintf("[%s:%d] ", file, line)
	return logrus.WithFields(fields)
}
