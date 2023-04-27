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
	WarnLogger := log.New(os.Stdout, "WARN ", goodFlag)

	WarnLogger.Println(v...)
}

func ErrorLogger(format string, v ...any) {
	ErrorLogger := log.New(os.Stderr, "ERROR ", badFlags)

	if format == "" {
		ErrorLogger.Println(v...)
	}

	ErrorLogger.Printf(format, v...)
}

func FatalLogger(v ...any) {
	FatalLogger := log.New(os.Stdout, "FATAL ", badFlags)

	FatalLogger.Fatal(v...)
}
