package logger

import (
	"log"
	"os"
)

var (
	logFlag = log.LstdFlags
)

func InfoLogger(format string, v ...any) {
	InfoLogger := log.New(os.Stdout, "INFO ", logFlag)

	if format == "" {
		InfoLogger.Println(v...)
	}

	InfoLogger.Printf(format, v...)
}

func WarnLogger(v ...any) {
	WarnLogger := log.New(os.Stdout, "WARN ", logFlag)

	WarnLogger.Println(v...)
}

func ErrorLogger(format string, v ...any) {
	ErrorLogger := log.New(os.Stderr, "ERROR ", logFlag)

	if format == "" {
		ErrorLogger.Println(v...)
	}

	ErrorLogger.Printf(format, v...)
}

func FatalLogger(v ...any) {
	FatalLogger := log.New(os.Stdout, "FATAL ", logFlag)

	FatalLogger.Fatal(v...)
}
