package util

func Filter[T comparable](in []T, predicate func(T) bool) []T {
	out := []T{}

	for _, item := range in {
		result := predicate(item)
		if result {
			out = append(out, item)
		}
	}
	return out
}
