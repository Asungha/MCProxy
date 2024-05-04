package utils

func SearchUTF8Byte(b []byte) string {
	var s string
	for _, v := range b {
		if v > 0x1F && v < 0x7F {
			s += string(v)
		} else {
			s += "."
		}
	}
	return s
}
