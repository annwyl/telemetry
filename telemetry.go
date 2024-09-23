package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
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
	transactions map[string]*Transaction
	mutex        sync.Mutex
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
	_, err := fmt.Fprintf(f.file, "%s %s %s %s\n", log.Timestamp.Format(time.RFC3339), log.Level, log.Message, log.Tags)
	return err
}

func (f *FileDriver) Close() error {
	return f.file.Close()
}

func NewLogger(driver Driver) *Logger {
	return &Logger{
		driver:       driver,
		transactions: make(map[string]*Transaction),
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

func (l *Logger) StartTransaction() string {
	transactionID := fmt.Sprintf("%d", time.Now().UnixNano()) // create a more unique id
	l.mutex.Lock()
	l.transactions[transactionID] = &Transaction{
		ID:    transactionID,
		Start: time.Now(),
	}
	l.mutex.Unlock()
	return transactionID
}

func (l *Logger) EndTransaction(transactionID string) error {
	l.mutex.Lock()
	transaction, exists := l.transactions[transactionID]
	if !exists {
		l.mutex.Unlock()
		return fmt.Errorf("transaction %s doesnt exist", transactionID)
	}
	transaction.End = time.Now()
	delete(l.transactions, transactionID)
	l.mutex.Unlock()
	// maybe log a summary or smth on how lnog it took etc, easier for kibana etc
	return nil
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

	jsonDriver, err := NewJSONDriver("log.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonDriver.Close()

	logger = NewLogger(jsonDriver)
	transID := logger.StartTransaction()
	logger.Debug("This is a debug message", map[string]string{"CPU": "CPU usage at 50%"}, transID)
	logger.EndTransaction(transID)
}
