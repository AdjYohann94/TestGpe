package service

import (
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/hashicorp/logutils"
)

const (
	LevelInfo        = "INFO"
	LevelDebug       = "DEBUG"
	LevelWarning     = "WARN"
	LevelError       = "ERROR"
	LogFilePath      = "log"
	LogFileDirectory = "log/debug.log"
)

// NewLogging initialize the logging interface to handle log levels.
// Different log levels are defined : INFO, DEBUG, WARN, ERROR.
//
// The minimum log level is defined in main.yml configuration file.
//
// Logs are saved into log/debug.log file.
func NewLogging() {
	_, err := os.Stat(LogFilePath)
	if os.IsNotExist(err) {
		err := os.Mkdir(LogFilePath, os.ModePerm)
		if err != nil {
			panic("cannot create log directory")
		}
	}

	f, err := os.OpenFile(LogFileDirectory, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic("cannot open debug.log file")
	}

	minLevel := viper.GetString(KeyLogLevel)
	log.SetOutput(&logutils.LevelFilter{
		Levels:   []logutils.LogLevel{LevelInfo, LevelDebug, LevelWarning, LevelError},
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   f,
	})
}
