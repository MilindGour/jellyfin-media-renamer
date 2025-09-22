package mediainfoprovider

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
)

type TmdbMIProvider struct {
	baseUrl  string
	scrapper scrapper.Scrapper
}

func NewTmdbMIProvider() *TmdbMIProvider {
	return &TmdbMIProvider{
		baseUrl:  "https://www.themoviedb.org",
		scrapper: scrapper.NewGoQueryScrapper(),
	}
}
func NewMockTmdbMIProvider() *TmdbMIProvider {
	return &TmdbMIProvider{
		baseUrl:  "tmdb",
		scrapper: scrapper.NewMockGoQueryScrapper(),
	}
}

func (t *TmdbMIProvider) SearchMediaInfo(term string, year int, mediaType MediaType) []MediaInfo {
	var url string
	var out []MediaInfo

	switch mediaType {
	case MediaTypeMovie:
		url = fmt.Sprintf("%s/search/movie?query=%s", t.baseUrl, t.getSearchString(term, year))
		out = t.getParsedMediaInfoListFromUrl(url, ".search_results.movie .card")
	case MediaTypeTV:
		url = fmt.Sprintf("%s/search/tv?query=%s", t.baseUrl, t.getSearchString(term, year))
		out = t.getParsedMediaInfoListFromUrl(url, ".search_results.tv .card")
	}

	return out
}

func (t *TmdbMIProvider) SearchMovies(term string, year int) []MovieResult {
	mediaInfoList := t.SearchMediaInfo(term, year, MediaTypeMovie)

	out := []MovieResult{}
	for _, info := range mediaInfoList {
		out = append(out, MovieResult{
			MediaInfo: info,
		})
	}
	return out
}

func (t *TmdbMIProvider) SearchTVShows(term string, year int) []TVResult {
	mediaInfoList := t.SearchMediaInfo(term, year, MediaTypeTV)

	out := []TVResult{}
	for _, info := range mediaInfoList {
		seasonList := t.getSeasonInformation(info.MediaID)
		out = append(out, TVResult{
			MediaInfo:    info,
			TotalSeasons: len(seasonList),
			Seasons:      seasonList,
		})
	}

	return out
}

func (t *TmdbMIProvider) GetJellyfinCompatibleDirectoryName(info MediaInfo) string {
	out := info.Name

	if info.YearOfRelease > 0 {
		out += fmt.Sprintf(" (%d)", info.YearOfRelease)
	}
	if len(info.MediaID) > 0 {
		out += fmt.Sprintf(" [tmdbid-%s]", info.MediaID)
	}

	return out
}

func (t *TmdbMIProvider) getParsedMediaInfoListFromUrl(url string, itemSelector string) []MediaInfo {
	scrapResult, err := t.scrapper.Scrap(url, itemSelector, t.getSearchItemFieldmap())
	if err != nil {
		return []MediaInfo{}
	}
	return t.parseScrapResultListToMediaInfo(scrapResult)
}

func (t *TmdbMIProvider) getSeasonInformation(mediaID string) []SeasonInfo {
	url := fmt.Sprintf("%s/tv/%s/seasons", t.baseUrl, mediaID)
	// seasonList := t.getParsedSeasonListFromUrl(url, ".media .column_wrapper .season_wrapper")
	seasonItemFieldmap := t.getSeasonItemFieldmap()
	scrapResult, err := t.scrapper.Scrap(url, ".media .column_wrapper .season_wrapper", seasonItemFieldmap)

	out := []SeasonInfo{}
	if err != nil {
		return out
	}
	for _, sr := range scrapResult {
		si := SeasonInfo{}

		// Season number
		seasonNumber, _ := sr["seasonNumber"]
		si.Number = t.extractSeasonNumber(seasonNumber)

		// Season total episodes
		seasonEpisodes, _ := sr["seasonTotalEpisodes"]
		si.TotalEpisodes = t.extractSeasonTotalEpisodes(seasonEpisodes)

		out = append(out, si)
	}

	return out
}
func (t *TmdbMIProvider) getSeasonItemFieldmap() map[string]string {
	return map[string]string{
		"seasonNumber":        ".panel .season a[href]",
		"seasonTotalEpisodes": ".panel .season .content h4",
	}
}

func (t *TmdbMIProvider) getSearchItemFieldmap() map[string]string {
	return map[string]string{
		"name":          ".title h2",
		"subname":       ".title h2 span.title",
		"description":   ".overview",
		"yearOfRelease": ".title .release_date",
		"thumbnailUrl":  ".image img.poster[src]",
		"mediaId":       ".title a.result[href]",
	}
}
func (t *TmdbMIProvider) parseScrapResultListToMediaInfo(in scrapper.ScrapResultList) []MediaInfo {
	out := []MediaInfo{}
	for _, scrapResult := range in {
		m := MediaInfo{}
		name, _ := scrapResult["name"]
		subName, _ := scrapResult["subname"]
		if len(subName) > 0 {
			fmt.Printf("Has subname: <%s>\n", subName)
			name = strings.ReplaceAll(name, subName, "")
		}
		m.Name = t.trimString(name)

		description, _ := scrapResult["description"]
		m.Description = t.trimString(description)

		yearOfRelease, _ := scrapResult["yearOfRelease"]
		m.YearOfRelease = t.extractYear(yearOfRelease)

		thumbnailUrl, _ := scrapResult["thumbnailUrl"]
		m.ThumbnailURL = t.trimString(thumbnailUrl)

		mediaId, _ := scrapResult["mediaId"]
		m.MediaID = t.extractMediaId(mediaId)

		out = append(out, m)
	}
	return out
}
func (t *TmdbMIProvider) trimString(in string) string {
	return strings.Trim(in, " \t\n\r")
}
func (t *TmdbMIProvider) extractYear(in string) int {
	// 20 September 2012
	out := 0
	yearRe := regexp.MustCompile(`(\d{4})`)
	match := yearRe.FindStringSubmatch(in)
	if match != nil {
		out, _ = strconv.Atoi(match[1])
	}
	return out
}
func (t *TmdbMIProvider) extractMediaId(in string) string {
	// /movie/12345-movie-name-2344
	out := ""
	if strings.Contains(in, "/") {
		lastSlashIndex := strings.LastIndex(in, "/")
		in = in[lastSlashIndex+1:]
	}
	idRe := regexp.MustCompile(`^(\d+)`)
	match := idRe.FindStringSubmatch(in)
	if match != nil {
		out = match[1]
	}

	return out
}
func (t *TmdbMIProvider) extractSeasonNumber(in string) int {
	// /xxx/season/2
	out := -1
	if strings.Contains(in, "/season/") {
		splits := strings.Split(in, "/season/")
		if len(splits) == 2 {
			str := splits[1]
			seasonNumber, err := strconv.Atoi(str)
			if err != nil {
				return out
			}
			out = seasonNumber
		}
	}

	return out
}

func (t *TmdbMIProvider) extractSeasonTotalEpisodes(in string) int {
	fmt.Println("In:", in)
	out := -1
	re := regexp.MustCompile(`(\d+) Episodes`)
	match := re.FindStringSubmatch(in)
	if match != nil {
		totalEpisodes, err := strconv.Atoi(match[1])
		if err != nil {
			return out
		}
		out = totalEpisodes
	}

	return out
}

func (t *TmdbMIProvider) getSearchString(name string, year int) string {
	out := url.QueryEscape(name)
	if year > 0 {
		out += fmt.Sprintf("%%20y:%d", year)
	}

	return out
}
