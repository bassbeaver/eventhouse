package logger

import (
	"github.com/bassbeaver/logopher"
	"os"
	"path/filepath"
)

const (
	LoggerFactoryServiceAlias = "LoggerFactory"
)

type LoggerFactory struct {
	logFilePath string
}

func (f *LoggerFactory) CreateLogger(bufferSize int) (*logopher.Logger, error) {
	if 0 == bufferSize {
		bufferSize = 1
	}

	logDir := filepath.Dir(f.logFilePath)
	_, dirStatError := os.Stat(logDir)
	if os.IsNotExist(dirStatError) {
		mkDirError := os.MkdirAll(logDir, 0644)
		if nil != mkDirError {
			panic(mkDirError)
		}
	} else if nil != dirStatError {
		panic(dirStatError)
	}

	logfile, logfileError := os.OpenFile(f.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if logfileError != nil {
		panic(logfileError)
	}

	fileHandler := logopher.CreateStreamHandler(logfile, &logopher.JsonFormatter{}, nil, nil, bufferSize)

	logger := &logopher.Logger{}
	logger.SetHandlers([]logopher.HandlerInterface{fileHandler})

	return logger, nil
}

//--------------------

func NewLoggerFactory(logFilePath string) *LoggerFactory {
	return &LoggerFactory{
		logFilePath: logFilePath,
	}
}
