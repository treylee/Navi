// utils/logger.go
package utils

import (
	"log"
	"os"
	"sync"
)

var (
	logInstance *log.Logger
	once        sync.Once
)

// GetLogger initializes and returns the global logger instance
func GetLogger() *log.Logger {
	once.Do(func() {
		// Initialize the logger only once
		logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		logInstance = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return logInstance
}

// LogError logs an error message
func LogError(message string, err error) {
	logger := GetLogger()
	logger.Printf("[ERROR] %s: %v\n", message, err)
}

// LogInfo logs an informational message
func LogInfo(message string) {
	logger := GetLogger()
	logger.Printf("[INFO] %s\n", message)
}
