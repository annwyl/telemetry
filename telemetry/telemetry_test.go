package telemetry

import (
	"encoding/json"
	"testing"
	"time"
)

type MockDriver struct {
	logs []Log
}

func (m *MockDriver) Log(log Log) error {
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockDriver) Close() error {
	return nil
}

func TestNewLogger(t *testing.T) {
	RegisterDriver("mock", func(config json.RawMessage) (Driver, error) {
		return &MockDriver{}, nil
	})

	config := Config{
		Name:        "mock",
		Config:      json.RawMessage(`{}`),
		LogLevel:    InfoLevel,
		DefaultTags: map[string]string{"environment": "test"},
	}

	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	if _, ok := logger.driver.(*MockDriver); !ok {
		t.Error("Logger isnt Mockdriver")
	}
}

func TestLogLevels(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := &Logger{
		driver: mockDriver,
		config: Config{
			LogLevel:    InfoLevel,
			DefaultTags: map[string]string{"environment": "test"},
		},
	}

	tests := []struct {
		level   LogLevel
		message string
		logFunc func(string, map[string]string, ...string) error
	}{
		{DebugLevel, "Debug message", logger.Debug},
		{InfoLevel, "Info message", logger.Info},
		{WarningLevel, "Warning message", logger.Warning},
		{ErrorLevel, "Error message", logger.Error},
	}

	for _, tt := range tests {
		err := tt.logFunc(tt.message, nil)
		if err != nil {
			t.Errorf("unexpected error for %v: %v", tt.level, err)
		}
	}

	if len(mockDriver.logs) != 4 {
		t.Errorf("wanted 4 logs, got %d", len(mockDriver.logs))
	}

	for i, log := range mockDriver.logs {
		if log.Level != tests[i].level {
			t.Errorf("wanted log level %v, got %v", tests[i].level, log.Level)
		}
		if log.Message != tests[i].message {
			t.Errorf("wanted message '%s', got '%s'", tests[i].message, log.Message)
		}
		if log.Tags["environment"] != "test" {
			t.Errorf("wanted default tag 'environment: test', got '%s'", log.Tags["environment"])
		}
	}
}

func TestTransaction(t *testing.T) {
	logger := &Logger{
		driver:       &MockDriver{},
		transactions: make(map[string]*Transaction),
	}

	transactionID := logger.StartTransaction()
	if transactionID == "" {
		t.Error("starttransaciton returned empty ID")
	}

	time.Sleep(100 * time.Millisecond)

	err := logger.EndTransaction(transactionID)
	if err != nil {
		t.Errorf("endtransaction returned error: %v", err)
	}

	if len(logger.transactions) != 0 {
		t.Error("transaction still exists")
	}

	err = logger.EndTransaction(transactionID)
	if err == nil {
		t.Error("endtransaction should return error as empty")
	}
}
