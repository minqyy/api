package str

func CompleteStringToLength(s string, length int, char rune) string {
	if length < len(s) {
		return s[:length]
	}
	buf := make([]rune, length-len(s))
	for i := 0; i < len(buf); i++ {
		buf[i] = char
	}
	return s + string(buf)
}
