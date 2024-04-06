package slice

func Contains[T comparable](s []T, el T) bool {
	for _, a := range s {
		if a == el {
			return true
		}
	}
	return false
}
