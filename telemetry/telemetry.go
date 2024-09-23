package telemetry

import (
	"fmt"
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

func NewLogger(config Config) *Logger {
	driver, err := getDriver(config)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &Logger{
		driver:       driver,
		transactions: make(map[string]*Transaction),
	}
}

func (l *Logger) Close() error {
	return l.driver.Close()
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
	transactionID := generateTransactionID()
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

func generateTransactionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // make it still more unique
}