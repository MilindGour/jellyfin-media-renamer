package renamer

import (
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	m "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
)

type JmrRenamerV1 struct {
	mip m.MediaInfoProvider
}

func NewJmrRenamerV1(mip m.MediaInfoProvider) *JmrRenamerV1 {
	return &JmrRenamerV1{
		mip: mip,
	}
}

func (j *JmrRenamerV1) GetMediaNameAndYear(rawFilename string) MediaNameAndYear {
	// Step 1: remove special characters from input
	onlyBasename := path.Base(rawFilename)
	outputFilename := removeSpecialCharacters(onlyBasename)

	// Step 2: extract year from input, if present otherwise 0
	// also remove all string after the year (including year)
	outputYear := extractYear(&outputFilename)

	// Step 3: remove all double spaces due to previous steps
	outputFilename = removeDoubleWhitespace(outputFilename)

	return MediaNameAndYear{outputFilename, outputYear}
}
func (j *JmrRenamerV1) GetMediaSeasonAndEpisode(rawFilepath string) MediaSeasonAndEpisode {
	filepathWithoutExtension, _ := strings.CutSuffix(rawFilepath, path.Ext(rawFilepath))
	in := strings.ToLower(filepathWithoutExtension)
	in += "_" // added to pass the last two regexps

	testREs := []*regexp.Regexp{
		regexp.MustCompile(`s(\d{2})e(\d{2})`),                     // SXXEXX
		regexp.MustCompile(`season[ \-_]+(\d{1,2})[ \-_]+(\d{2})`), // Season X - XX
		regexp.MustCompile(`s(\d+)[ \-_]+(\d+)`),                   // SX - XX
		regexp.MustCompile(`episode[ \-_]+(\d+)`),                  // Episode XX (No season information)
		regexp.MustCompile(`[^0-9](\d{1,2})[^0-9]`),                // XX (No season information)
	}

	for _, testre := range testREs {
		m1 := testre.FindStringSubmatch(in)
		if m1 != nil {
			if len(m1) == 2 {
				// contains only episode number
				episode, e2 := strconv.Atoi(m1[1])
				if e2 != nil {
					panic("Cannot parse episode number. " + e2.Error())
				}
				if episode > 200 {
					return MediaSeasonAndEpisode{-1, -1}
				}
				return MediaSeasonAndEpisode{1, episode}
			} else if len(m1) == 3 {
				// contains both episode and season number
				season, e1 := strconv.Atoi(m1[1])
				if e1 != nil {
					panic("Cannot parse season number. " + e1.Error())
				}
				episode, e2 := strconv.Atoi(m1[2])
				if e2 != nil {
					panic("Cannot parse episode number. " + e2.Error())
				}
				return MediaSeasonAndEpisode{season, episode}
			}
		}
	}

	return MediaSeasonAndEpisode{-1, -1}
}

func removeSpecialCharacters(inputFilename string) string {
	outputFilename := ""
	for _, ch := range inputFilename {
		re := regexp.MustCompile("[A-Za-z0-9 ]")

		if re.MatchString(string(ch)) {
			outputFilename += string(ch)
		} else {
			outputFilename += " "
		}
	}
	return outputFilename
}

func extractYear(filename *string) int {
	re := regexp.MustCompile(" ?[0-9]{4,} ?")

	if re.MatchString(*filename) {
		matchedYearString := re.FindString(*filename)
		matchedYearIndex := re.FindStringIndex(*filename)
		if len(matchedYearString) == 0 {
			return -1
		}
		yearInt, err := strconv.Atoi(strings.Trim(matchedYearString, " "))
		if err != nil {
			log.Fatal("Cannot convert year string to integer", err)
		}
		*filename = (*filename)[0:matchedYearIndex[0]]

		if yearInt > 9999 || yearInt < 1900 {
			yearInt = -1
		}
		return int(yearInt)
	}
	return -1
}

func removeDoubleWhitespace(str string) string {
	re := regexp.MustCompile(" {2,}")

	singleSpacedStr := re.ReplaceAllString(str, " ")
	return strings.Trim(singleSpacedStr, " ")
}
