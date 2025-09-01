package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	parser := NewSteakParser(config)

	// Run parsing immediately on startup
	log.Println("Running initial steak parsing...")
	err = parser.ParseAndNotify()
	if err != nil {
		log.Printf("Error during initial parsing: %v", err)
	}

	c := cron.New()
	_, err = c.AddFunc(config.Tracking.Interval, func() {
		log.Println("Starting steak parsing...")
		err := parser.ParseAndNotify()
		if err != nil {
			log.Printf("Error during parsing: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	fmt.Println("Steak parser started. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	c.Stop()
}
