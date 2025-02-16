package util

func P(s string) *string {
	return &s
}

func Clamp(a, low, high int) int {
	return min(max(a, low), high)
}

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
