package scrapper_test

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
)

func TestTmdbScrapper_GetSearchableString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   models.ClearFileEntry
		want string
	}{
		{
			name: "Get search string with year",
			in:   models.ClearFileEntry{Name: "Test", Year: 2012},
			want: "Test y:2012",
		},
		{
			name: "Get search string without year",
			in:   models.ClearFileEntry{Name: "Test", Year: 0},
			want: "Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := scrapper.NewTmdbScrapper()
			got := tm.GetSearchableString(tt.in)
			if got != tt.want {
				t.Errorf("GetSearchableString() = %v, want %v", got, tt.want)
			}
		})
	}
}
