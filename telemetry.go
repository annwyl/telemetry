package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	driver       Driver
	transacitons map[string]*Transaction
}

type Transaction struct {
	ID    string
	Start time.Time
	End   time.Time
	Logs  []Log
}

type Log struct {
	Timestamp     time.Time
	Level         LogLevel
	Message       string
	Tags          map[string]string
	TransactionID string
}

type Driver interface {
	Log(log Log) error
	Close() error
}

type ConsoleDriver struct{}

type JSONDriver struct {
	file    *os.File
	encoder *json.Encoder
}

type FileDriver struct {
	file *os.File
}

func (c *ConsoleDriver) Log(log Log) error {
	fmt.Println(log)
	return nil
}

func NewJSONDriver(filename string) (*JSONDriver, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &JSONDriver{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (j *JSONDriver) Log(log Log) error {
	return j.encoder.Encode(log)
}

func (j *JSONDriver) Close() error {
	return j.file.Close()
}

func NewFileDriver(filename string) (*FileDriver, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileDriver{file: file}, nil
}

func (f *FileDriver) Log(log Log) error {
	_, err := fmt.Fprintf(f.file, "%s %s %s %s\n", log.Timestamp, log.Level, log.Message, log.Tags)
	return err
}

func (f *FileDriver) Close() error {
	return f.file.Close()
}

func NewLogger(driver Driver) *Logger {
	return &Logger{
		driver:       driver,
		transacitons: make(map[string]*Transaction),
	}
}

func (l *Logger) log(level LogLevel, message string, tags map[string]string, transactionID string) error {
	log := Log{
		Timestamp:     time.Now(),
		Level:         level,
		Message:       message,
		Tags:          tags,
		TransactionID: transactionID,
	}
	return l.driver.Log(log)
}

func (l *Logger) Debug(message string, tags map[string]string, transactionID string) error {
	return l.log(DebugLevel, message, tags, transactionID)
}

func (l *Logger) Info(message string, tags map[string]string, transactionID string) error {
	return l.log(DebugLevel, message, tags, transactionID)
}

func (l *Logger) Warning(message string, tags map[string]string, transactionID string) error {
	return l.log(DebugLevel, message, tags, transactionID)
}

func (l *Logger) Error(message string, tags map[string]string, transactionID string) error {
	return l.log(DebugLevel, message, tags, transactionID)
}

func main() {
	/*logger := NewLogger()
	logger.Debug("This is a debug message", map[string]string{"CPU": "CPU usage at 50%"})
	logger.Info("This is an info message", map[string]string{"CPU": "CPU usage at 60%"})
	logger.Warning("This is a warning message", map[string]string{"CPU": "CPU usage at 90%"})
	logger.Error("This is an error message", map[string]string{"CPU": "CPU usage at 100%"})*/

	fileDriver, err := NewFileDriver("log.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileDriver.Close()

	logger := NewLogger(fileDriver)
	logger.Debug("This is a debug message", map[string]string{"CPU": "CPU usage at 50%"}, "123")
}
