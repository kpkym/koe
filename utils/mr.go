package utils

func Filter[T any](items []T, fn func(item T) bool) []T {
	result := make([]T, 0)

	for _, e := range items {
		if fn(e) {
			result = append(result, e)
		}
	}
	return result
}

func Map[S, T any](items []S, fn func(item S) T) []T {
	result := make([]T, 0)

	for _, e := range items {
		result = append(result, fn(e))
	}
	return result

}
