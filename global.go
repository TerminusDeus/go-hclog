package hclog

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	lumberjack "github.com/TerminusDeus/lumberjack"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

var (
	AgentOptions []*LoggerOptions

	protect sync.Once
	def     Logger

	// DefaultOptions is used to create the Default logger. These are read
	// only when the Default logger is created, so set them as soon as the
	// process starts.
	DefaultOptions = &LoggerOptions{
		Level:  DefaultLevel,
		Output: DefaultOutput,
		TimeFn: time.Now,
	}
)

// Default returns a globally held logger. This can be a good starting
// place, and then you can use .With() and .Named() to create sub-loggers
// to be used in more specific contexts.
// The value of the Default logger can be set via SetDefault() or by
// changing the options in DefaultOptions.
//
// This method is goroutine safe, returning a global from memory, but
// care should be used if SetDefault() is called it random times
// in the program as that may result in race conditions and an unexpected
// Logger being returned.
func Default() Logger {
	protect.Do(func() {
		// If SetDefault was used before Default() was called, we need to
		// detect that here.
		if def == nil {
			def = New(DefaultOptions)
		}
	})

	return def
}

// L is a short alias for Default().
func L() Logger {
	return Default()
}

// SetDefault changes the logger to be returned by Default()and L()
// to the one given. This allows packages to use the default logger
// and have higher level packages change it to match the execution
// environment. It returns any old default if there is one.
//
// NOTE: This is expected to be called early in the program to setup
// a default logger. As such, it does not attempt to make itself
// not racy with regard to the value of the default logger. Ergo
// if it is called in goroutines, you may experience race conditions
// with other goroutines retrieving the default logger. Basically,
// don't do that.
func SetDefault(log Logger) Logger {
	old := def
	def = log
	return old
}

func SetAgentOptions(options []*LoggerOptions) {
	AgentOptions := make([]*LoggerOptions, 0, len(options))

	// assumes that several destinations are set
	for _, opts := range options {
		prepareOptions(opts)

		AgentOptions = append(AgentOptions, opts)
	}
}

func prepareOptions(opts *LoggerOptions) {
	if opts.LogFile != "-" {
		opts.JSONFormat = opts.LogFormat == "json"
		opts.Level = LevelFromString(opts.LogLevel)

		if opts.LogPath != "" {
			logFileName := opts.LogPath + "/"

			if opts.LogFile != "" {
				logFileName += opts.LogFile
			} else {
				logFileName += fmt.Sprintf("new_log_file_%s", time.Now().String())
			}

			f, err := os.Create(logFileName)
			if err != nil {
				panic(err)
			}

			f.Close()

			logFileMaxSizeRaw := opts.LogMaxSize

			var logFileMaxSize int
			if logFileMaxSizeRaw != "" {

				size, err := parseutil.ParseCapacityString(logFileMaxSizeRaw)
				if err != nil {
					panic(errors.New("bad value for log_max_size: " + logFileMaxSizeRaw))
				}

				logFileMaxSize = int(size)
			}

			logFileTTLRaw := opts.LogRotate

			var logFileTTL int

			if logFileTTLRaw != "" {
				dur, err := parseutil.ParseDurationSecond(logFileTTLRaw)
				if err != nil {
					panic(errors.New("bad value for log_rotate: " + logFileTTLRaw))
				}

				logFileTTL = int(dur.Seconds())
			}

			opts.Output = &lumberjack.Logger{
				Filename: logFileName,
				MaxSize:  logFileMaxSize, // bytes
				MaxAge:   logFileTTL,     // seconds
			}
		}
	}
}
