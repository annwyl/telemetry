package telemetry

import (
	"crypto/rand"
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
	config       Config
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
		config:       config,
		transactions: make(map[string]*Transaction),
	}
}

func (l *Logger) Close() error {
	return l.driver.Close()
}

func (l *Logger) log(level LogLevel, message string, tags map[string]string, transactionID ...string) error {
	if tags == nil {
		tags = make(map[string]string)
	}

	for k, v := range l.config.DefaultTags {
		tags[k] = v
	}

	log := Log{
		Timestamp:     time.Now(),
		Level:         level,
		Message:       message,
		Tags:          tags,
		TransactionID: "",
	}

	if len(transactionID) == 1 {
		log.TransactionID = transactionID[0]
	}

	return l.driver.Log(log)
}

func (l *Logger) Debug(message string, tags map[string]string, transactionID ...string) error {
	return l.log(DebugLevel, message, tags, transactionID...)
}

func (l *Logger) Info(message string, tags map[string]string, transactionID ...string) error {
	return l.log(InfoLevel, message, tags, transactionID...)
}

func (l *Logger) Warning(message string, tags map[string]string, transactionID ...string) error {
	return l.log(WarningLevel, message, tags, transactionID...)
}

func (l *Logger) Error(message string, tags map[string]string, transactionID ...string) error {
	return l.log(ErrorLevel, message, tags, transactionID...)
}

func (l *Logger) SetLogLevel(level LogLevel) {
	l.config.LogLevel = level
}

func (l *Logger) AddDefaultTag(key, value string) {
	l.config.DefaultTags[key] = value
}

func (l *Logger) DeleteDefaultTag(key string) {
	delete(l.config.DefaultTags, key)
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
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}
