package regnet

func subStr(s string, pos, length int) string {
	bytes := []byte(s)
	l := pos + length
	if l > len(bytes) {
		l = len(bytes)
	}
	return string(bytes[pos:l])
}
