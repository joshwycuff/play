package model

var index = -1

type Result struct {
	ExitStatus int
	Stdout     string
	Stderr     string
}

type HistoryEntry struct {
	command string
	result  Result
}

var history = []HistoryEntry{}

func PushHistoryEntry(entry HistoryEntry) int {
	if len(history) > 0 && entry == history[len(history)-1] {
		return index
	}
	history = append(history, entry)
	index = len(history) - 1
	return index
}

func GetPreviousHistoryEntry() HistoryEntry {
	if len(history) == 0 {
		return HistoryEntry{}
	}
	index = clamp(index-1, 0, len(history)-1)
	return history[index]
}

func GetHistoryEntry() HistoryEntry {
	if len(history) == 0 {
		return HistoryEntry{}
	}
	return history[clamp(index, 0, len(history)-1)]
}

func GetNextHistoryEntry() HistoryEntry {
	if len(history) == 0 {
		return HistoryEntry{}
	}
	index = clamp(index+1, 0, len(history)-1)
	return history[index]
}

func clamp(a, low, high int) int {
	return min(max(a, low), high)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
