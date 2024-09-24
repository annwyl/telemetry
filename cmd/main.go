package main

import (
	"fmt"
	"os"

	_ "github.com/annwyl/telemetry/drivers"
	"github.com/annwyl/telemetry/telemetry"
)

func main() {
	config, err := telemetry.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("failed to load config")
		os.Exit(1)
	}

	fmt.Println("Registered drivers:")
	for name := range telemetry.GetRegisteredDrivers() {
		fmt.Println(name)
	}

	logger := telemetry.NewLogger(config)
	if logger == nil {
		fmt.Println("failed to create logger")
		os.Exit(1)
	}
	defer logger.Close()

	transactionID := logger.StartTransaction()

	if transactionID == "" {
		fmt.Println("failed to start transaction")
		os.Exit(1)
	}

	err = logger.Debug("Thhis is a debug message", map[string]string{"CPU": "CPU usage is at 69%"}, transactionID)
	if err != nil {
		fmt.Println("failed debug message")
		os.Exit(1)
	}

	err = logger.Info("This is a info mesage", map[string]string{"CPU": "CPU usage is at 69%"}, transactionID)
	if err != nil {
		fmt.Println("Failed info message")
		os.Exit(1)
	}

	err = logger.Warning("This is a warning message", map[string]string{"CPU": "CPU usage is at 69%"})
	if err != nil {
		fmt.Println("failed warning message")
		os.Exit(1)
	}

	err = logger.EndTransaction(transactionID)
	if err != nil {
		fmt.Println("failed to end transaction")
		os.Exit(1)
	}

}
