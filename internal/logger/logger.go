package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	accessLog  *log.Logger
)

func Init(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", path, err)
	}

	infoLogger = log.New(f, "INFO ", log.LstdFlags)
	errorLogger = log.New(f, "ERROR ", log.LstdFlags|log.Lshortfile)
	accessLog = log.New(f, "", 0)

	return nil
}

func Info(format string, v ...any) {
	if infoLogger != nil {
		infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...any) {
	if errorLogger != nil {
		errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Access(format string, v ...any) {
	if accessLog != nil {
		accessLog.Output(2, fmt.Sprintf(format, v...))
	}
}