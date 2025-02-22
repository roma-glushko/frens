package utils

func Unique[T comparable](s []T) []T {
	m := make(map[T]struct{}, len(s))

	for _, v := range s {
		m[v] = struct{}{}
	}

	unique := make([]T, 0, len(m))

	for k := range m {
		unique = append(unique, k)
	}

	return unique
}
