package logger

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	SuccessLevel
	ProgressLevel
	ProgressErrorLevel
	ProgressDebugLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type CustomLogger struct {
	*log.Logger
}

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "SUCCESS", "PROGRESS", "PROGRESS_ERROR", "PROGRESS_DEBUG", "WARN", "ERROR", "FATAL"}[l]
}

func NewCustomLogger(packageName string) *CustomLogger {
	logger := log.New(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	if packageName == "" {
		logger.SetStyles(GetDefaultLoggerStyle())
	} else {
		logger.SetStyles(GetPackageLoggerStyle(packageName))
	}

	return &CustomLogger{Logger: logger}
}

var DefaultLogger = NewCustomLogger("main")

func (c *CustomLogger) Log(message *MessageWrapper) {
	switch message.LogLevel {
	case DebugLevel:
		c.Debug(message.Message)
	case InfoLevel:
		c.Info(message.Message)
	case SuccessLevel:
		c.Success(message.Message)
	case ProgressLevel:
		c.Progress(message.Message)
	case WarnLevel:
		c.Warn(message.Message)
	case ErrorLevel:
		c.Error(message.Message)
	case FatalLevel:
		c.Fatal(message.Message)
	default:
		c.Info(message.Message)
	}
}

func Initialize(verbose bool) {
	DefaultLogger.Log(MsgInitializing)

	if verbose {
		DefaultLogger.Debug("running in verbose mode")
		DefaultLogger.SetLevel(log.DebugLevel)
	}

	log.SetStyles(GetDefaultLoggerStyle())
}

func (c *CustomLogger) Success(msg interface{}, keyvals ...interface{}) {
	c.Logger.Log(log.Level(SuccessLevel), msg, keyvals...)
}

func (c *CustomLogger) Successf(format string, args ...interface{}) {
	c.Logger.Log(log.Level(SuccessLevel), fmt.Sprintf(format, args...))
}

func (c *CustomLogger) Progress(msg interface{}, keyvals ...interface{}) {
	c.Logger.Log(log.Level(ProgressLevel), msg, keyvals...)
}

func (c *CustomLogger) Progressf(format string, args ...interface{}) {
	c.Logger.Log(log.Level(ProgressLevel), fmt.Sprintf(format, args...))
}

func (c *CustomLogger) ProgressError(msg interface{}, keyvals ...interface{}) {
	c.Logger.Log(log.Level(ProgressErrorLevel), msg, keyvals...)
}

func (c *CustomLogger) ProgressErrorf(format string, args ...interface{}) {
	c.Logger.Log(log.Level(ProgressErrorLevel), fmt.Sprintf(format, args...))
}

func (c *CustomLogger) ProgressDebug(msg interface{}, keyvals ...interface{}) {
	c.Logger.Log(log.Level(ProgressDebugLevel), msg, keyvals...)
}

func (c *CustomLogger) ProgressDebugf(format string, args ...interface{}) {
	c.Logger.Log(log.Level(ProgressDebugLevel), fmt.Sprintf(format, args...))
}