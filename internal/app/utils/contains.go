package utils

func Contains[T comparable](s []T, sub T) bool {
	for _, v := range s {
		if sub == v {
			return true
		}
	}
	return false
}
