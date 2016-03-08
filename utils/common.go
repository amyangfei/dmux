package utils

// In returns true if target is in candidates, otherwise returns false
func In(target string, candidates []string) bool {
	for _, item := range candidates {
		if target == item {
			return true
		}
	}
	return false
}
