package utils

func Contains(slice []string, elem string) bool {
	for _, s := range slice {
		if s == elem {
			return true
		}
	}
	return false
}

func AllEmpty(strSlice []string) bool {
	for _, s := range strSlice {
		if len(s) > 0 {
			return false
		}
	}
	return true
}
