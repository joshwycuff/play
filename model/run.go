package model

import (
	"bytes"
	"os/exec"
	"syscall"
)

func run(command string, stdin string) (HistoryEntry, error) {
	historyEntry := HistoryEntry{command: command, result: Result{}}

	cmd := exec.Command("sh", "-c", command)

	// Get the command's stdin pipe
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return historyEntry, err
	}

	// Provide input to the command
	go func() {
		defer stdinPipe.Close()
		stdinPipe.Write([]byte(stdin))
	}()

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			status := exitErr.Sys().(syscall.WaitStatus)
			historyEntry.result.ExitStatus = status.ExitStatus()
			historyEntry.result.Stdout = stdout.String()
			historyEntry.result.Stderr = stderr.String()
		} else {
			return historyEntry, err
		}
		// m.output.Failure()
		// m.output.SetContent(err.Error())
	} else {
		historyEntry.result.ExitStatus = 0
		historyEntry.result.Stdout = stdout.String()
		historyEntry.result.Stderr = stderr.String()

		// m.output.Success()
		// m.output.SetContent(stdout.String())
	}

	return historyEntry, nil
}
