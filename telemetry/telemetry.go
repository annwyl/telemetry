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

func NewLogger(config Config) (*Logger, error) {
	driver, err := getDriver(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	return &Logger{
		driver:       driver,
		config:       config,
		transactions: make(map[string]*Transaction),
	}, nil
}

func (l *Logger) Close() error {
	return l.driver.Close()
}

func (l *Logger) log(level LogLevel, message string, tags map[string]string, transactionID ...string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if level < l.config.LogLevel {
		return nil
	}

	count := len(tags) + len(l.config.DefaultTags)

	var finalTags map[string]string
	if count > 0 {
		finalTags = make(map[string]string, count)

		for k, v := range l.config.DefaultTags {
			finalTags[k] = v
		}

		for k, v := range tags {
			finalTags[k] = v
		}
	}

	log := Log{
		Timestamp:     time.Now(),
		Level:         level,
		Message:       message,
		Tags:          finalTags,
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
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.config.LogLevel = level
}

func (l *Logger) AddDefaultTag(key, value string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.config.DefaultTags[key] = value
}

func (l *Logger) DeleteDefaultTag(key string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.config.DefaultTags, key)
}

func (l *Logger) StartTransaction() string {
	transactionID := generateTransactionID()
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.transactions[transactionID] = &Transaction{
		ID:    transactionID,
		Start: time.Now(),
	}
	return transactionID
}

func (l *Logger) EndTransaction(transactionID string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	transaction, exists := l.transactions[transactionID]
	if !exists {
		return fmt.Errorf("endtransaction %s doesnt exist", transactionID)
	}
	transaction.End = time.Now()
	delete(l.transactions, transactionID)
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
