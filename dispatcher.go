package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Dispatch(log *os.File, job JobConfig) {
	timestamp := func() string { return time.Now().Format("2006-01-02 15:04:05") }

	fmt.Fprintf(log, "[%s] - %s START\n", timestamp(), job.Name)

	if strings.HasPrefix(job.Command, "logrotate") {
		parts := strings.Fields(job.Command)
		if len(parts) == 3 {
			LogRotateStart(parts[1], parts[2])
		}
		fmt.Fprintf(log, "[%s] - %s END\n", timestamp(), job.Name)
	} else {
		cmd := exec.Command("powershell.exe", job.Command)
		output, err := cmd.CombinedOutput()
		log.Write(output)
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
		}
		fmt.Fprintf(log, "[%s] - %s END exitcode=%d\n", timestamp(), job.Name, exitCode)
	}
}
