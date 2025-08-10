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
		scrapper: scrapper.NewGoQuery(),
	}
}
func NewMockTmdbMIProvider() *TmdbMIProvider {
	return &TmdbMIProvider{
		baseUrl:  "tmdb",
		scrapper: scrapper.NewMockGoQuery(),
	}
}

func (t *TmdbMIProvider) SearchMovies(term string) []MovieResult {
	url := fmt.Sprintf("%s/search/movie?query=%s", t.baseUrl, url.PathEscape(term))
	mediaInfoList := t.getParsedMediaInfoListFromUrl(url, ".search_results.movie .card")

	out := []MovieResult{}
	for _, info := range mediaInfoList {
		out = append(out, MovieResult{
			MediaInfo: info,
		})
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

func (t *TmdbMIProvider) SearchTVShows(term string) []TVResult {
	// TODO: Implement this method
	panic("Not implemented")
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
		m.MediaID = t.extraMediaId(mediaId)

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
func (t *TmdbMIProvider) extraMediaId(in string) string {
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
