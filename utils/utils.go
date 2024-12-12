package utils

// Go inexplicably does not have an abs function for ints, and I have
// been bitten by float casting far too many times to do it that way,
// so we write our own.
func Distance(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}
