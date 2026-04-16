package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println("process-cron launcher")
		fmt.Println("schedule in task manager, configure by schedule.yaml")
		fmt.Println()
		fmt.Println("runs every n minutes (30) recommended and executes scripts")
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	dayName := now.Format("Mon")

	data, err := os.ReadFile("schedule.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read schedule.yaml: %v\n", err)
		os.Exit(1)
	}

	var config ConfigRoot
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse schedule.yaml: %v\n", err)
		os.Exit(1)
	}

	statePath := "state.json"
	state := make(map[string]string)
	if stateData, err := os.ReadFile(statePath); err == nil {
		json.Unmarshal(stateData, &state)
	}

	for _, job := range config.Jobs {
		if len(job.Days) > 0 && !contains(job.Days, dayName) {
			continue
		}

		runAt, err := time.Parse("15:04", job.RunAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid run_at for job %s: %v\n", job.Name, err)
			continue
		}

		runToday := time.Date(now.Year(), now.Month(), now.Day(), runAt.Hour(), runAt.Minute(), 0, 0, now.Location())

		if last, ok := state[job.Name]; ok && last == today {
			continue
		}
		if now.Before(runToday) {
			continue
		}

		os.MkdirAll("locks", 0755)
		lockFile := fmt.Sprintf("locks/%s.lock", job.Name)
		if _, err := os.Stat(lockFile); err == nil {
			continue
		}

		os.WriteFile(lockFile, []byte(runToday.Format("02.01.2006 15:04:05.0000")), 0644)
		os.MkdirAll("logs", 0755)
		logPath := fmt.Sprintf("logs/%s.log", job.Name)

		func() {
			defer os.Remove(lockFile)

			logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open log %s: %v\n", logPath, err)
				return
			}
			defer logFile.Close()

			Dispatch(logFile, job)
			state[job.Name] = today
		}()
	}

	stateJSON, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(statePath, stateJSON, 0644)
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}
