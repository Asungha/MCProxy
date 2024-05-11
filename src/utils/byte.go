package utils

func Concat(slice ...[]byte) []byte {
	var result []byte
	for _, s := range slice {
		result = append(result, s...)
	}
	return result
}
