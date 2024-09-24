package telemetry

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"
)

type MockDriver struct {
	logs []Log
	mu   sync.Mutex
}

func (m *MockDriver) Log(log Log) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockDriver) Close() error {
	return nil
}

func TestNewLogger(t *testing.T) {
	err := RegisterDriver("mockNewLogger", func(config json.RawMessage) (Driver, error) {
		return &MockDriver{}, nil
	})

	if err != nil {
		t.Fatalf("registerdriver gave error: %v", err)
	}

	config := Config{
		Name:        "mockNewLogger",
		Config:      json.RawMessage(`{}`),
		LogLevel:    InfoLevel,
		DefaultTags: map[string]string{"environment": "test"},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("newlogger returned error: %v", err)
	}

	if _, ok := logger.driver.(*MockDriver); !ok {
		t.Error("Logger isn't Mockdriver")
	}
}

func TestLogLevels(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := &Logger{
		driver: mockDriver,
		config: Config{
			LogLevel:    DebugLevel,
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

func TestSetLogLevel(t *testing.T) {
	logger := &Logger{
		driver: &MockDriver{},
		config: Config{LogLevel: InfoLevel},
	}

	logger.SetLogLevel(DebugLevel)
	if logger.config.LogLevel != DebugLevel {
		t.Errorf("wanted log level to be debuglevel, got %v", logger.config.LogLevel)
	}
}

func TestDefaultTags(t *testing.T) {
	logger := &Logger{
		driver: &MockDriver{},
		config: Config{DefaultTags: make(map[string]string)},
	}

	logger.AddDefaultTag("app_name", "telemetry")
	if logger.config.DefaultTags["app_name"] != "telemetry" {
		t.Errorf("wanted default tag 'key: value', got '%s'", logger.config.DefaultTags["key"])
	}

	logger.DeleteDefaultTag("app_name")
	if _, exists := logger.config.DefaultTags["app_name"]; exists {
		t.Error("wanted default tag deleted, still exists")
	}
}

func TestUniqueTransactions(t *testing.T) {
	logger := &Logger{
		driver:       &MockDriver{},
		transactions: make(map[string]*Transaction),
	}

	trx1 := logger.StartTransaction()
	trx2 := logger.StartTransaction()
	if trx1 == trx2 {
		t.Error("transaction should be unique")
	}
}

func TestConcurrentLogging(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := &Logger{
		driver: mockDriver,
		config: Config{
			LogLevel:    InfoLevel,
			DefaultTags: map[string]string{"environment": "test"},
		},
	}

	var wg sync.WaitGroup
	logCount := 10

	wg.Add(logCount)
	for i := 0; i < logCount; i++ {
		go func(i int) {
			defer wg.Done()
			err := logger.Info("threaded log message", map[string]string{"environment": "test"})
			if err != nil {
				t.Errorf("error in threaded logging: %v", err)
			}
		}(i)
	}

	wg.Wait()

	if len(mockDriver.logs) != logCount {
		t.Errorf("wanted %d logs, have %d", logCount, len(mockDriver.logs))
	}
}

func TestConfigurationOverride(t *testing.T) {
	err := RegisterDriver("mockConfigOverride", func(config json.RawMessage) (Driver, error) {
		return &MockDriver{}, nil
	})

	if err != nil {
		t.Fatalf("registerdriver gave error: %v", err)
	}

	configContent := `{
		"driver": "mockConfigOverride",
		"driver_config": "",
		"log_level": 1,
		"default_tags": {"environment": "test"}
	}`
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}

	config.LogLevel = ErrorLevel
	config.DefaultTags["environment"] = "production"

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("newlogger returned error: %v", err)
	}

	if logger.config.LogLevel != ErrorLevel {
		t.Errorf("wanted loglevel %v, got %v", ErrorLevel, logger.config.LogLevel)
	}

	if logger.config.DefaultTags["environment"] != "production" {
		t.Errorf("wanted default tag 'environment: production', got '%s'", logger.config.DefaultTags["environment"])
	}
}

// probably also test if timestamps are correct, lots of logs, long transactions
