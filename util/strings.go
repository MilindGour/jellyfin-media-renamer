package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// This function cleans the filenames to make it searchable in scrapped sites
func CleanFilename(inputFilename string) (string, int) {
	// "American Sniper (2014)",

	// Step 1: remove special characters from input
	var outputFilename string = removeSpecialCharacters(inputFilename)

	// Step 2: extract year from input, if present otherwise 0
	// also remove all string after the year (including year)
	var outputYear int = extractYear(&outputFilename)

	// Step 3: remove all double spaces due to previous steps
	outputFilename = removeDoubleWhitespace(outputFilename)

	return outputFilename, outputYear
}

func JoinPaths(paths ...string) string {
	out := ""
	for _, path := range paths {
		trimmedPath := strings.TrimRight(path, "/")
		if len(out) > 0 {
			out += "/"
		}
		out += trimmedPath
	}
	return out
}

func ExtractYearFromString(in string) (int, error) {
	yearRe := regexp.MustCompile(`(\d{4})`)
	match := yearRe.FindString(in)
	return strconv.Atoi(match)
}

func removeSpecialCharacters(inputFilename string) string {
	outputFilename := ""
	for _, ch := range inputFilename {
		re, err := regexp.Compile("[A-Za-z0-9 ]")
		if err != nil {
			log.Fatal("Unable to compile regular expression")
		}
		if re.MatchString(string(ch)) {
			outputFilename += string(ch)
		} else {
			outputFilename += " "
		}
	}
	return outputFilename
}

func extractYear(filename *string) int {
	re, err := regexp.Compile(" ?[0-9]{4,} ?")
	if err != nil {
		log.Fatal("Unable to compile regular expression")
	}
	if re.MatchString(*filename) {
		matchedYearString := re.FindString(*filename)
		matchedYearIndex := re.FindStringIndex(*filename)
		if len(matchedYearString) == 0 {
			return 0
		}
		yearInt, err := strconv.Atoi(strings.Trim(matchedYearString, " "))
		if err != nil {
			log.Fatal("Cannot convert year string to integer", err)
		}
		*filename = (*filename)[0:matchedYearIndex[0]]

		if yearInt > 9999 {
			yearInt = 0
		}
		return int(yearInt)
	}
	return 0
}

func removeDoubleWhitespace(str string) string {
	re, err := regexp.Compile(" {2,}")
	if err != nil {
		log.Fatal("Cannot compile regex:", err)
	}
	singleSpacedStr := re.ReplaceAllString(str, " ")
	return strings.Trim(singleSpacedStr, " ")
}
