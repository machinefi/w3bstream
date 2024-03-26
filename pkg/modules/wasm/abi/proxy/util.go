package proxy

// subStringWithLength
// If the length is negative, an empty string is returned.
// If the length is greater than the length of the input string, the entire string is returned.
// Otherwise, a substring of the input string with the specified length is returned.
func subStringWithLength(str string, length int) string {
	if length < 0 {
		return ""
	}
	rs := []rune(str)
	strLen := len(rs)

	if length > strLen {
		return str
	}
	return string(rs[0:length])
}
