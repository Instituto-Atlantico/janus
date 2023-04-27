package logger

import (
	"log"
	"os"
)

var (
	goodFlag = log.LstdFlags
	badFlags = log.LstdFlags | log.Lshortfile
)

func InfoLogger(format string, v ...any) {
	InfoLogger := log.New(os.Stdout, "INFO ", goodFlag)

	if format == "" {
		InfoLogger.Println(v...)
	}

	InfoLogger.Printf(format, v...)
}

func WarnLogger(v ...any) {
	logger := log.New(os.Stdout, "WARN ", goodFlag)

	logger.Println(v...)
}

func ErrorLogger(v ...any) {
	logger := log.New(os.Stdout, "ERROR ", badFlags)

	logger.Println(v...)
}

func FatalLogger(v ...any) {
	logger := log.New(os.Stdout, "FATAL ", badFlags)

	logger.Fatal(v...)
}
