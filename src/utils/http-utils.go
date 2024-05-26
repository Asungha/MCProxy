package utils

func IsHTTPMethod(s string) bool {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, a := range methods {
		if a == s {
			return true
		}
	}
	return false
}
