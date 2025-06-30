package tests

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/util"
)

var dirtyFilenames []string = []string{
	"American Sniper (2014)",
	"Captain America - Civil War",
	"Contracted.Phase.II.2015.BRRip.XviD-ETRG",
	"Deadpool (2016) 720p BluRay x264 [Dual Audio] [Hindi (Line Audio) - English] ESubs - Downloadhub",
	"Fifty Shades of Grey (2015)",
	"Harry Potter & The Chamber of Secrets 2002 BRRip 720p Dual-Audio[Eng+Hindi] ~ BRAR",
	"Harry Potter & The Deathly Hallows Part I 2010 BRRip 720p Dual-Audio[Eng+Hindi] ~ BRAR",
	"Housefull 3 2016 Hindi (1CD) DvDScr x264 AAC - Hon3y",
	"Random movie 12342",
}

var cleanFilenames []string = []string{
	"American Sniper",
	"Captain America Civil War",
	"Contracted Phase II",
	"Deadpool",
	"Fifty Shades of Grey",
	"Harry Potter The Chamber of Secrets",
	"Harry Potter The Deathly Hallows Part I",
	"Housefull 3",
	"Random movie",
}

var cleanYears []int = []int{
	2014,
	0,
	2015,
	2016,
	2015,
	2002,
	2010,
	2016,
	0,
}

func TestCleanFilename(t *testing.T) {
	for index, input := range dirtyFilenames {
		wantFilename := cleanFilenames[index]
		wantYear := cleanYears[index]
		actualFilename, actualYear := util.CleanFilename(input)
		if actualFilename != wantFilename {
			t.Errorf("\nWanted name: %s\nActual name: %s", wantFilename, actualFilename)
		}
		if actualYear != wantYear {
			t.Errorf("\nWanted year: %d\nActual year: %d", wantYear, actualYear)
		}
	}
}

func TestExtractMediaIdFromUrl(t *testing.T) {
	inputs := []string{"/tv/1234-test-tv", "https://abs.com/tv/123-sdf", "/movie/1231"}
	expectedOut := []string{"1234", "123", "1231"}

	for i := range inputs {
		in := inputs[i]
		want := expectedOut[i]
		actual := util.ExtractMediaIdFromUrl(in)

		if want != actual {
			t.Errorf("\nInput: %s\nWanted Id: %s\nActual Id: %s", in, want, actual)
		}
	}
}

func TestExtractTotalEpisodesFromInfoString(t *testing.T) {
	// 1990 • 27 Episodes
	input := "1990 • 27 Episodes"
	want := 27
	actual := util.ExtractTotalEpisodesFromInfoString(input)

	if want != actual {
		t.Errorf("\nInput: %s\nWanted Number: %d\nActual Number: %d", input, want, actual)
	}
}
