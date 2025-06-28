package util

func Filter[T any](in []T, predicate func(T) bool) []T {
	out := []T{}

	for _, item := range in {
		result := predicate(item)
		if result {
			out = append(out, item)
		}
	}
	return out
}

// Environment related code
type Environment int

const (
	DEV  Environment = 0
	PROD Environment = 1
)

var env Environment = PROD

func SetEnvironment(e Environment) {
	env = e
}
func IsProduction() bool {
	return env == PROD
}
