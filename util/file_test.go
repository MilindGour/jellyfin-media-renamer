package util_test

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func TestSortBySeasonAndEpisodeNumbers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		a    models.MediaPathRename
		b    models.MediaPathRename
		want int
	}{
		{
			name: "sort result when a > b on season",
			a:    models.MediaPathRename{NewPath: "S02E01"},
			b:    models.MediaPathRename{NewPath: "S01E05"},
			want: 1,
		},
		{
			name: "sort result when a > b on episode",
			a:    models.MediaPathRename{NewPath: "S01E02"},
			b:    models.MediaPathRename{NewPath: "S01E01"},
			want: 1,
		},
		{
			name: "sort result when b > a on season",
			a:    models.MediaPathRename{NewPath: "S01E02"},
			b:    models.MediaPathRename{NewPath: "S03E01"},
			want: -1,
		},
		{
			name: "sort result when b > a on episode",
			a:    models.MediaPathRename{NewPath: "S01E02"},
			b:    models.MediaPathRename{NewPath: "S01E05"},
			want: -1,
		},
		{
			name: "sort result when a == b",
			a:    models.MediaPathRename{NewPath: "S01E02"},
			b:    models.MediaPathRename{NewPath: "S01E02"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.SortBySeasonAndEpisodeNumbers(tt.a, tt.b)
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("SortBySeasonAndEpisodeNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
