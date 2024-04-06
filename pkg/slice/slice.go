package slice

// Contains check if slice contains an element
func Contains[T comparable](s []T, el T) bool {
	for _, a := range s {
		if a == el {
			return true
		}
	}
	return false
}
