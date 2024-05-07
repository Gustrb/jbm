package utils

func Contains(arr []string, val string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			return true
		}
	}

	return false
}
