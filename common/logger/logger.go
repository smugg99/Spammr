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
	WarnLevel
	ErrorLevel
	FatalLevel
)

type CustomLogger struct {
	*log.Logger
}

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "SUCCESS", "WARN", "ERROR", "FATAL"}[l]
}

func (c *CustomLogger) Success(msg interface{}, keyvals ...interface{}) {
	c.Logger.Log(log.Level(SuccessLevel), msg, keyvals...)
}

func (c *CustomLogger) Successf(format string, args ...interface{}) {
	c.Logger.Log(log.Level(SuccessLevel), fmt.Sprintf(format, args...))
}

func NewCustomLogger(packageName string) *CustomLogger {
	logger := log.New(os.Stdout)

	logger.SetLevel(log.DebugLevel)
	logger.SetStyles(GetPackageLoggerStyle(packageName))

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
