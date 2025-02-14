package util

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	// Open a log file for writing
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not open debug log file:", err)
		os.Exit(1)
	}

	// Create a logger
	logger = log.New(f, "DEBUG: ", log.LstdFlags)
}

func Debug(msg string) {
	logger.Println(msg)
}
