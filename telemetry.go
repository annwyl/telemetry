package main

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

type Logger struct {
}

type Log struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(level LogLevel, message string) Log {
	return Log{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
}

func (l *Logger) Debug(message string) {
	logEntry := l.log(DebugLevel, message)
	fmt.Println(logEntry)
}

func (l *Logger) Info(message string) {
	logEntry := l.log(InfoLevel, message)
	fmt.Println(logEntry)
}

func (l *Logger) Warning(message string) {
	logEntry := l.log(WarningLevel, message)
	fmt.Println(logEntry)
}

func (l *Logger) Error(message string) {
	logEntry := l.log(ErrorLevel, message)
	fmt.Println(logEntry)
}

func main() {
	logger := NewLogger()
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warning("This is a warning message")
	logger.Error("This is an error message")
}
