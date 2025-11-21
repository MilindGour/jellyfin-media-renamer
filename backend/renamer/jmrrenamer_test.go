package renamer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
	"github.com/MilindGour/jellyfin-media-renamer/testdata"
	"github.com/MilindGour/jellyfin-media-renamer/websocket"
)

func TestJmrRenamerV1_GetMediaNameAndYear(t *testing.T) {
	type fields struct {
		mip mediainfoprovider.MediaInfoProvider
	}
	type args struct {
		rawFilename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MediaNameAndYear
	}{
		{
			"Simple movie name",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"Test movie name [2023]"},
			MediaNameAndYear{"Test movie name", 2023},
		},
		{
			"Intermediate movie name",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"Test.movie name- [2025] New TMRIP"},
			MediaNameAndYear{"Test movie name", 2025},
		},
		{
			"Advanced movie name",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"Airplane! (1980) [BluRay] [1080p] [YTS.LT]"},
			MediaNameAndYear{"Airplane", 1980},
		},
		{
			"Movie name without year",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"Test Movie 2"},
			MediaNameAndYear{"Test Movie 2", -1},
		},
		{
			"Movie name with invalid year",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"Special movie 1080p"},
			MediaNameAndYear{"Special movie", -1},
		},
		{
			"Blank movie name",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{""},
			MediaNameAndYear{"", -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JmrRenamer{
				mip: tt.fields.mip,
			}
			if got := j.GetMediaNameAndYear(tt.args.rawFilename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JmrRenamerV1.GetMediaNameAndYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJmrRenamerV1_GetMediaSeasonAndEpisode(t *testing.T) {
	type fields struct {
		mediainfoprovider.MediaInfoProvider
	}
	type args struct {
		rawFilename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MediaSeasonAndEpisode
	}{
		{
			"Simple show name",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"Simple Show Name"},
			MediaSeasonAndEpisode{-1, -1},
		},
		{
			"Show With episode only",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"some/show/Show Name 05"},
			MediaSeasonAndEpisode{1, 5},
		},
		{
			"Simple season and episode",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"some/good/show file s01e02"},
			MediaSeasonAndEpisode{1, 2},
		},
		{
			"Show name variant 2",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"some/good/show file season 3 - 04"},
			MediaSeasonAndEpisode{3, 4},
		},
		{
			"Show name variant 3",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"some/good/show file s3 - 12"},
			MediaSeasonAndEpisode{3, 12},
		},
		{
			"Show with [2x03] in name",
			fields{mediainfoprovider.NewMockTmdbMIProvider()},
			args{"some/good/show file [2x03]"},
			MediaSeasonAndEpisode{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JmrRenamer{
				mip: tt.fields.MediaInfoProvider,
			}
			if got := j.GetMediaSeasonAndEpisode(tt.args.rawFilename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JmrRenamerV1.ExtractTVShowFileInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJmrRenamer_ConfirmEntriesForRename(t *testing.T) {
	var request RenameMediaConfirmRequest
	err := json.Unmarshal(testdata.RenameConfirmRequestMock, &request)
	if err != nil {
		t.Errorf("ConfirmRename() err: %s", err.Error())
	}

	fmt.Printf("Total entries: %d\n", len(request))

	j := NewJmrRenamerV1(mediainfoprovider.NewMockTmdbMIProvider(), filesystem.NewJmrFS(), config.NewDevJmrConfig(), websocket.NewJMRWebSocket())
	j.ConfirmEntriesForRename(request)
}
