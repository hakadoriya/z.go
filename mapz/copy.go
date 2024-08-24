package mapz

func Copy[T any](m map[string]T) map[string]T {
	n := make(map[string]T, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}
