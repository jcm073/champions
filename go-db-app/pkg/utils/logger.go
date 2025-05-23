package utils

import (
	"log"
	"os"
)

var (
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
	infoLogger.Println(message)
}

func Error(message string) {
	errorLogger.Println(message)
}

func Debug(message string) {
	debugLogger.Println(message)
}