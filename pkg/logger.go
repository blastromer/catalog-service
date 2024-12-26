package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger struct holds the log file and the log level
type Logger struct {
	file   *os.File
	logger *log.Logger
}

// NewLogger initializes a new logger instance
func NewLogger(logFile string) (*Logger, error) {
	// Open the log file for appending, create it if it doesn't exist
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Create a new logger with a specific log format
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Return the new logger instance
	return &Logger{file: file, logger: logger}, nil
}

// Info logs an informational message
func (l *Logger) Info(message string) {
	l.logger.SetPrefix("[INFO] ")
	l.logger.Println(message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.logger.SetPrefix("[ERROR] ")
	l.logger.Println(message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.logger.SetPrefix("[DEBUG] ")
	l.logger.Println(message)
}

// Close closes the log file
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("Error closing log file:", err)
	}
}

// LogToConsole writes logs to both the console and a file
func (l *Logger) LogToConsole(message string, level string) {
	logMessage := fmt.Sprintf("%s: %s", level, message)
	fmt.Println(logMessage)
	l.logger.Println(logMessage)
}

func main() {
	// Initialize the logger
	logFile := "app.log"
	logger, err := NewLogger(logFile)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	defer logger.Close()

	// Example usage
	logger.Info("This is an info message")
	logger.Debug("This is a debug message")
	logger.Error("This is an error message")

	// Log to console as well
	logger.LogToConsole("This is a message to both console and file", "[INFO]")
}
