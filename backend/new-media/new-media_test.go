package newmedia

import (
	"strings"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/network"
)

func TestNewMedia_SearchMedia(t *testing.T) {
	type args struct {
		searchTerm string
	}
	tests := []struct {
		name             string
		args             args
		wantTotalResults int
	}{
		{
			name: "Check ability to parse response",
			args: args{
				searchTerm: "TestSearchTerm",
			},
			wantTotalResults: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewMedia{
				h: network.NewMockResponse(),
			}
			if got := n.SearchMedia(tt.args.searchTerm); len(got) != tt.wantTotalResults {
				t.Errorf("NewMedia.SearchMedia() = %d results, want %d results", len(got), tt.wantTotalResults)
			}
		})
	}
}

func TestNewMedia_GetMagneticURL(t *testing.T) {
	tests := []struct {
		name      string
		item      NewMediaSearchItem
		wantStart string
	}{
		{
			name: "get mag url",
			item: NewMediaSearchItem{
				ID:       "100",
				Name:     "Test 1",
				InfoHash: "1234567890",
			},
			wantStart: "magnet:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNewMedia(network.NewMockResponse())
			got := n.GetMagneticURL(tt.item)
			if strings.Index(got, tt.wantStart) > 0 {
				t.Errorf("GetMagneticURL() should start with magnet:")
			}
		})
	}
}
