package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/byitkc/gosumo"
)

// Log is a representation of a basic log message.
type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func main() {
	logs := generateLogs()
	sumoClient := gosumo.Client{
		Endpoint: "<endpoint url>",
	}
	if err := gosumo.PostLogs(sumoClient, logs); err != nil {
		slog.Error("error posting logs to Sumo Logic", "error", err)
		os.Exit(1)
	}
}

// generateLogs generates a couple of dummy logs for use with development
func generateLogs() []Log {
	var logs []Log
	log1 := Log{
		Timestamp: time.Now(),
		Message:   "This is a test log #1",
	}
	time.Sleep(time.Second * 1)
	log2 := Log{
		Timestamp: time.Now(),
		Message:   "This is a test log #2",
	}
	logs = append(logs, log1)
	logs = append(logs, log2)
	return logs
}
