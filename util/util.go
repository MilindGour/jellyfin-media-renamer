package util

func Filter[T any](in []T, predicateFn func(T) bool) []T {
	out := []T{}
	for _, t := range in {
		if predicateFn(t) {
			out = append(out, t)
		}
	}
	return out
}
func Find[T any](in []T, predicateFn func(T) bool) *T {
	for _, t := range in {
		if predicateFn(t) {
			return &t
		}
	}
	return nil
}
