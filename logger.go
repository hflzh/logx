// Package logx provides logging functionality for Go applications.
package logx

import (
	"fmt"
	"io"
	"log"
)

// Level represents a logging level.
//
// A set of standard logging levels are defined as follows:
//
//   Fine
//   Debug
//   Info
//   Warn
//   Error
//
// These levels are ordered, from lowest priority to highest priority.
// Enabling logging at a given level also enables logging at all higher levels.
// For example, if the desired logging level for a logger is `Debug`, the
// messages of Debug level, as well as Info, Warn and Error levels are all logged.
//
// There is a special level `Off` that can be used to turn off logging.
type Level int

const (
	// Fine is used to log fine granularity messages.
	Fine Level = 10

	// Debug is used to log debug messages.
	Debug Level = 20

	// Info is used to log informative messages.
	Info Level = 30

	// Warn is used to log warning messages.
	Warn Level = 40

	// Error is used to log error messages.
	Error Level = 50

	// Off turns off logging.
	Off Level = 99
)

// String returns a string representation for the logging level.
//
// This implements the fmt.Stringer interface.
func (level Level) String() string {
	switch level {
	case Fine:
		return "Fine"
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warn:
		return "Warn"
	case Error:
		return "Error"
	case Off:
		return "Off"
	default:
		return "Unknown"
	}
}

// Logger is a wrapper for log.Logger of the standard log package, adding
// capabilities to control the output of messges at desired level and whether
// the log time is displayed in local time zone or UTC.
type Logger struct {
	// The underlying logger.
	logger *log.Logger

	// The desired logging level for the logger.
	level Level

	// The suffix of log time.
	// This is "UTC " if using UTC time, or is an empty string if using local time.
	timeSuffix string
}

// New creates a logger that writes messages of the logging level to io.Writer.
//
// If the specified level is Off or does not match any of the standard logging
// levels, it returns nil representing logging is disabled.
func New(out io.Writer, level Level, useLocalTime bool) *Logger {
	if out == nil {
		return nil
	}

	switch level {
	case Fine, Debug, Info, Warn, Error:
	case Off:
		return nil
	default:
		return nil
	}

	var tz string
	flag := log.LstdFlags | log.Lmicroseconds
	if !useLocalTime {
		flag |= log.LUTC
		tz = "UTC "
	}

	return &Logger{
		level:      level,
		logger:     log.New(out, "", flag),
		timeSuffix: tz,
	}
}

// Fine writes the specified message at Fine level to the output.
func (l *Logger) Fine(messageFormat string, messageArgs ...interface{}) {
	l.Log(Fine, messageFormat, messageArgs...)
}

// Debug writes the specified message at Debug level to the output.
func (l *Logger) Debug(messageFormat string, messageArgs ...interface{}) {
	l.Log(Debug, messageFormat, messageArgs...)
}

// Info writes the specified message at Info level to the output.
func (l *Logger) Info(messageFormat string, messageArgs ...interface{}) {
	l.Log(Info, messageFormat, messageArgs...)
}

// Warn writes the specified message at Warn level to the output.
func (l *Logger) Warn(messageFormat string, messageArgs ...interface{}) {
	l.Log(Warn, messageFormat, messageArgs...)
}

// Error writes the specified message at Error level to the output.
func (l *Logger) Error(messageFormat string, messageArgs ...interface{}) {
	l.Log(Error, messageFormat, messageArgs...)
}

// Log writes the specified message at the specified logging level to the output.
func (l *Logger) Log(level Level, messageFormat string, messageArgs ...interface{}) {
	if level == Off || l == nil || l.level > level {
		return
	}

	l.logger.Print(l.timeSuffix+label(level), fmt.Sprintf(messageFormat, messageArgs...))
}

// LogWithFn checks if the message at the specified logging level should be logged.
// If so, it calls the fn and writes the message returned by fn to the output.
func (l *Logger) LogWithFn(level Level, fn func() string) {
	if level == Off || l.level > level {
		return
	}

	l.logger.Print(l.timeSuffix+label(level), fn())
}

// label returns a label for the specified logging level.
func label(level Level) string {
	switch level {
	case Fine:
		return "[FINE]  "
	case Debug:
		return "[DEBUG] "
	case Info:
		return "[INFO]  "
	case Warn:
		return "[WARN]  "
	case Error:
		return "[ERROR] "
	default:
		return ""
	}
}
