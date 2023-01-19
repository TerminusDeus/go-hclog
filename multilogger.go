package hclog

import (
	"io"
	"log"
	"time"
)

// Make sure that multiLogger is a Logger
var _ Logger = &multiLogger{}

// multiLogger is an internal logger agent based implementation. Internal in that it is
// defined entirely by this package.
type multiLogger struct {
	nestedLoggers []*intLogger
}

type intLoggerPointer *intLogger

func newMultiLogger(multipleOpts []*LoggerOptions) *multiLogger {
	l := &multiLogger{nestedLoggers: make([]*intLogger, 0, len(multipleOpts))}
	for _, opts := range multipleOpts {
		l.nestedLoggers = append(l.nestedLoggers, newLogger(opts))
	}

	return l
}

// Log a message and a set of key/value pairs if the given level is at
// or more severe that the threshold configured in the Logger.
func (l *multiLogger) log(name string, level Level, msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.log(name, level, msg, args)
	}
}

// logPlain is the non-JSON logging format function which writes directly
// to the underlying writer the logger was initialized with.
//
// If the logger was initialized with a color function, it also handles
// applying the color to the log message.
//
// Color Options
//  1. No color.
//  2. Color the whole log line, based on the level.
//  3. Color only the header (level) part of the log line.
//  4. Color both the header and fields of the log line.
func (l *multiLogger) logPlain(t time.Time, name string, level Level, msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.logPlain(t, name, level, msg, args)
	}
}

// JSON logging function
func (l *multiLogger) logJSON(t time.Time, name string, level Level, msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.logJSON(t, name, level, msg, args)
	}
}

func (l multiLogger) jsonMapEntry(t time.Time, name string, level Level, msg string) map[string]interface{} {
	m := map[string]interface{}{}
	for _, logger := range l.nestedLoggers {
		nestedM := logger.jsonMapEntry(t, name, level, msg)
		for k, v := range nestedM {
			m[k] = v
		}
	}

	return m
}

// Emit the message and args at the provided level
func (l *multiLogger) Log(level Level, msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Log(level, msg, args)
	}
}

// Emit the message and args at DEBUG level
func (l *multiLogger) Debug(msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Debug(msg, args)
	}
}

// Emit the message and args at TRACE level
func (l *multiLogger) Trace(msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Trace(msg, args)
	}
}

// Emit the message and args at INFO level
func (l *multiLogger) Info(msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Info(msg, args)
	}
}

// Emit the message and args at WARN level
func (l *multiLogger) Warn(msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Warn(msg, args)
	}
}

// Emit the message and args at ERROR level
func (l *multiLogger) Error(msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Error(msg, args)
	}
}

// Indicate that the logger would emit TRACE level logs
func (l *multiLogger) IsTrace() bool {
	return false
}

// Indicate that the logger would emit DEBUG level logs
func (l *multiLogger) IsDebug() bool {
	return false
}

// Indicate that the logger would emit INFO level logs
func (l *multiLogger) IsInfo() bool {
	return false
}

// Indicate that the logger would emit WARN level logs
func (l *multiLogger) IsWarn() bool {
	return false
}

// Indicate that the logger would emit ERROR level logs
func (l *multiLogger) IsError() bool {
	return false
}

// Return a sub-Logger for which every emitted log message will contain
// the given key/value pairs. This is used to create a context specific
// Logger.
func (l *multiLogger) With(args ...interface{}) Logger {
	var newLoggers []*intLogger

	for _, logger := range l.nestedLoggers {
		newLoggers = append(l.nestedLoggers, logger.With(args...).(*intLogger))
	}

	l.nestedLoggers = nil
	l.nestedLoggers = newLoggers

	return l
}

// Create a new sub-Logger that a name decending from the current name.
// This is used to create a subsystem specific Logger.
func (l *multiLogger) Named(name string) Logger {
	var newLoggers []*intLogger

	for _, logger := range l.nestedLoggers {
		newLoggers = append(l.nestedLoggers, logger.Named(name).(*intLogger))
	}

	l.nestedLoggers = nil
	l.nestedLoggers = newLoggers

	return l
}

// Create a new sub-Logger with an explicit name. This ignores the current
// name. This is used to create a standalone logger that doesn't fall
// within the normal hierarchy.
func (l *multiLogger) ResetNamed(name string) Logger {
	var newLoggers []*intLogger

	for _, logger := range l.nestedLoggers {
		newLoggers = append(l.nestedLoggers, logger.Named(name).(*intLogger))
	}

	l.nestedLoggers = nil
	l.nestedLoggers = newLoggers

	return l
}

func (l *multiLogger) ResetOutput(opts *LoggerOptions) error {
	for _, logger := range l.nestedLoggers {
		logger.ResetOutput(opts)
	}

	return nil
}

func (l *multiLogger) ResetOutputWithFlush(opts *LoggerOptions, flushable Flushable) error {
	for _, logger := range l.nestedLoggers {
		logger.ResetOutputWithFlush(opts, flushable)
	}

	return nil
}

// Update the logging level on-the-fly. This will affect all subloggers as
// well.
func (l *multiLogger) SetLevel(level Level) {
	for _, logger := range l.nestedLoggers {
		logger.SetLevel(level)
	}
}

// Create a *log.Logger that will send it's data through this Logger. This
// allows packages that expect to be using the standard library log to actually
// use this logger.
func (l *multiLogger) StandardLogger(opts *StandardLoggerOptions) *log.Logger {
	// var newLoggers []*intLogger

	// for _, logger := range l.nestedLoggers {
	// 	newLoggers = append(l.nestedLoggers, logger.StandardLogger(opts).(*intLogger))
	// }

	// l.nestedLoggers = nil
	// l.nestedLoggers = newLoggers
	return nil
}

func (l *multiLogger) StandardWriter(opts *StandardLoggerOptions) io.Writer {
	return nil
}

// Accept implements the SinkAdapter interface
func (l *multiLogger) Accept(name string, level Level, msg string, args ...interface{}) {
	for _, logger := range l.nestedLoggers {
		logger.Accept(name, level, msg, args)
	}
}

// ImpliedArgs returns the loggers implied args
func (i *multiLogger) ImpliedArgs() []interface{} {
	return nil
}

// Name returns the loggers name
func (l *multiLogger) Name() string {
	for _, logger := range l.nestedLoggers {
		return logger.Name()
	}
	return ""
}
