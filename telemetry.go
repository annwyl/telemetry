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
	driver Driver
}

type Log struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Tags      map[string]string
}

type Driver interface {
	Log(log Log) error
}

type ConsoleDriver struct{}

func (c *ConsoleDriver) Log(log Log) error {
	fmt.Println(log)
	return nil
}

func NewLogger() *Logger {
	return &Logger{
		driver: &ConsoleDriver{},
	}
}

func (l *Logger) log(level LogLevel, message string, tags map[string]string) error {
	log := Log{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Tags:      tags,
	}
	return l.driver.Log(log)
}

func (l *Logger) Debug(message string, tags map[string]string) error {
	return l.log(DebugLevel, message, tags)
}

func (l *Logger) Info(message string, tags map[string]string) error {
	return l.log(DebugLevel, message, tags)
}

func (l *Logger) Warning(message string, tags map[string]string) error {
	return l.log(DebugLevel, message, tags)
}

func (l *Logger) Error(message string, tags map[string]string) error {
	return l.log(DebugLevel, message, tags)
}

func main() {
	logger := NewLogger()
	logger.Debug("This is a debug message", map[string]string{"CPU": "CPU usage at 50%"})
	logger.Info("This is an info message", map[string]string{"CPU": "CPU usage at 60%"})
	logger.Warning("This is a warning message", map[string]string{"CPU": "CPU usage at 90%"})
	logger.Error("This is an error message", map[string]string{"CPU": "CPU usage at 100%"})
}
