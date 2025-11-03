package main

// hasHeadlessFlag checks if the --headless flag is present in the provided arguments
func hasHeadlessFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--headless" {
			return true
		}
	}
	return false
}