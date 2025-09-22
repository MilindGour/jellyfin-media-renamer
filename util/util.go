package util

import (
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
)

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

func SortBySizeDesc(a, b filesystem.DirEntry) int {
	if a.Size > b.Size {
		return -1
	}
	if a.Size < b.Size {
		return 1
	}
	return 0
}
