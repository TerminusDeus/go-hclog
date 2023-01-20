package hclog

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	lumberjack "github.com/TerminusDeus/lumberjack"
)

var (
	agentOptions []*LoggerOptions

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
	agentOptions := make([]*LoggerOptions, 0, len(options))

	for _, opts := range options {
		if opts.LogFile != "" {
			logFileName := opts.LogPath + opts.LogFile

			fmt.Printf("opts.LogFile = %v\n", opts.LogFile)
			fmt.Printf("opts.LogMaxSize = %v\n", opts.LogMaxSize)
			fmt.Printf("opts.LogRotate = %v\n", opts.LogRotate)
			fmt.Printf("opts.LogFormat = %v\n", opts.LogFormat)
			fmt.Printf("opts.LogPath = %v\n", opts.LogPath)

			if _, err := os.Stat(logFileName); err == nil {
				_, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}

				logFileMaxSizeRaw := opts.LogMaxSize // os.Getenv("VAULT_AGENT_LOG_FILE_MAX_SIZE")

				var logFileMaxSize int
				if logFileMaxSizeRaw != "" {
					logFileMaxSize, err = strconv.Atoi(logFileMaxSizeRaw)
					if err != nil {
						panic(errors.New("bad value for logFileMaxSize: " + logFileMaxSizeRaw))
					}
				}

				logFileTTLRaw := opts.LogRotate // os.Getenv("VAULT_AGENT_LOG_FILE_MAX_AGE")

				var logFileTTL int
				if logFileTTLRaw != "" {
					logFileTTL, err = strconv.Atoi(logFileTTLRaw)
					if err != nil {
						panic(errors.New("bad value for logFileTTL: " + logFileTTLRaw))
					}
				}

				opts.JSONFormat = opts.LogFormat == "json"
				opts.Level = LevelFromString(opts.LogLevel)

				opts.Output = &lumberjack.Logger{
					Filename: logFileName,
					MaxSize:  logFileMaxSize, // megabytes
					MaxAge:   logFileTTL,     //minutes
				}
			}
			fmt.Printf("New: opts: %+v", opts)
		}
		agentOptions = append(agentOptions, opts)
	}
}

// type VaultAgentOptions struct {
// 	// vault agent specific options
// 	LogRotate  string
// 	LogMaxSize string
// 	LogFile    string
// 	LogPath    string
// 	LogFormat  string
// 	LogLevel   string
// }
