package model

import (
	"bytes"
	"os/exec"
	"syscall"

	"github.com/joshwycuff/play/util"
)

func run(command string, stdin *string) (HistoryData, error) {
	historyEntry := HistoryData{Command: command, Stdin: stdin}

	cmd := exec.Command("sh", "-c", command)

	// Get the command's stdin pipe
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return historyEntry, err
	}

	// Provide input to the command
	go func() {
		defer stdinPipe.Close()
		stdinPipe.Write([]byte(*stdin))
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
			historyEntry.Stdout = util.P(stdout.String())
			historyEntry.Stderr = util.P(stderr.String())
			status := exitErr.Sys().(syscall.WaitStatus)
			historyEntry.ExitStatus = status.ExitStatus()
		} else {
			return historyEntry, err
		}
	} else {
		historyEntry.Stdout = util.P(stdout.String())
		historyEntry.Stderr = util.P(stderr.String())
		historyEntry.ExitStatus = 0
	}

	return historyEntry, nil
}
