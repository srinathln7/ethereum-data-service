package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// ANSI escape codes for colors
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[96m"
	Reset = "\033[0m"
)

// CustomLogger wraps the standard logger to include colors
type CustomLogger struct {
	logger *log.Logger
}

// New creates a new instance of CustomLogger
func New(prefix string) *CustomLogger {
	return &CustomLogger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
	}
}

// Println logs a message in the specified color
func (c *CustomLogger) Println(color string, v ...interface{}) {
	c.logger.Println(color + currentTime() + fmt.Sprint(v...) + Reset)
}

// Printf logs a formatted message in the specified color
func (c *CustomLogger) Printf(color string, format string, v ...interface{}) {
	c.logger.Printf(color+currentTime()+format+Reset, v...)
}

// Success logs a success message in green
func (c *CustomLogger) Successln(v ...interface{}) {
	c.Println(Green, v...)
}

// Successf logs a formatted success message in green
func (c *CustomLogger) Successf(format string, v ...interface{}) {
	c.Printf(Green, format, v...)
}

// Info logs an informational message in blue
func (c *CustomLogger) Infoln(v ...interface{}) {
	c.Println(Blue, v...)
}

// Infof logs a formatted informational message in blue
func (c *CustomLogger) Infof(format string, v ...interface{}) {
	c.Printf(Blue, format, v...)
}

// Fatal logs a fatal message in red and exits
func (c *CustomLogger) Fatal(v ...interface{}) {
	c.Println(Red, v...)
	os.Exit(1)
}

// Fatalf logs a formatted fatal message in red and exits
func (c *CustomLogger) Fatalf(format string, v ...interface{}) {
	c.Printf(Red, format, v...)
	os.Exit(1)
}

// Utility function to get current time formatted with log.LstdFlags
func currentTime() string {
	return time.Now().Format(time.Stamp)
}
