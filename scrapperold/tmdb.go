package scrapperold

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gocolly/colly"
)

type TmdbScrapper struct {
	BaseURL string
}

func NewTmdbScrapper() *TmdbScrapper {
	out := TmdbScrapper{
		BaseURL: "https://www.themoviedb.org",
	}
	return &out
}

func (t TmdbScrapper) GetSearchableString(in models.ClearFileEntry) string {
	if in.Year == 0 {
		return in.Name
	}

	return fmt.Sprintf("%s y:%d", in.Name, in.Year)
}

func (t TmdbScrapper) SearchTV(in models.ClearFileEntry) ([]models.TVResult, error) {
	searchString := t.GetSearchableString(in)
	out := []models.TVResult{}

	log.Printf("Searching for tv: %s\n", searchString)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Accept-Language", "en-US")
	})

	c.OnHTML(".search_results.tv div.card.v4", func(h *colly.HTMLElement) {
		mediaInfo := t.ScrapMediaInfoListFromCollyElement(h)
		log.Println("MediaInfo:", mediaInfo)
		seasons := t.ScrapSeasonInfoList(mediaInfo.MediaID)

		result := models.TVResult{
			MediaInfo:    mediaInfo,
			TotalSeasons: len(seasons),
			Seasons:      seasons,
		}
		out = append(out, result)
	})

	// prepare for visiting the url
	pathEscape := url.PathEscape("query=" + searchString)
	visitURL := fmt.Sprintf("%s/search/movie?%s", t.BaseURL, pathEscape)

	err := c.Visit(visitURL)
	c.Wait()

	if err != nil {
		return nil, errors.New("Cannot visit search movie url. " + err.Error())
	}

	log.Println("Total results found for tv", in.Name, ":", len(out))
	return out, nil
}

func (t TmdbScrapper) ScrapMediaInfoListFromCollyElement(h *colly.HTMLElement) models.MediaInfo {
	nameNodeString := h.ChildText(".title h2")
	nameTitleNodeString := h.ChildText(".title h2 .title")
	if len(nameTitleNodeString) > 0 {
		nameNodeString = strings.Replace(nameNodeString, nameTitleNodeString, "", 1)
	}
	nameNodeString = strings.Trim(nameNodeString, " ")
	mediaLink := h.ChildAttr(".details .title a.result", "href")
	mediaID := util.ExtractMediaIDFromURL(mediaLink)
	releaseDateStr := h.ChildText(".release_date")
	releaseYear, err := util.ExtractYearFromString(releaseDateStr)
	if err != nil {
		releaseYear = 0
	}

	result := models.MediaInfo{
		Name:          nameNodeString,
		Description:   h.ChildText(".overview"),
		ThumbnailURL:  h.ChildAttr(".poster img", "src"),
		MediaID:       mediaID,
		YearOfRelease: releaseYear,
	}
	return result
}

func (t TmdbScrapper) ScrapSeasonInfoList(mediaID string) []models.SeasonInfo {
	out := []models.SeasonInfo{}
	cc := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"),
	)

	cc.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Accept-Language", "en-US")
	})
	cc.OnHTML(".media.tv_v4 .season_wrapper", func(h *colly.HTMLElement) {
		seasonLink := h.ChildAttr(".season .content h2 a", "href")
		seasonNumberString := util.ExtractMediaIDFromURL(seasonLink)
		seasonNumber, err := strconv.Atoi(seasonNumberString)
		if err != nil {
			log.Println("Cannot extract season number. Defaulting to -1")
			seasonNumber = -1
		}
		seasonInfoString := h.ChildText(".season .content h4")
		seasonTotalEpisodes := util.ExtractTotalEpisodesFromInfoString(seasonInfoString)

		out = append(out, models.SeasonInfo{
			Number:        seasonNumber,
			TotalEpisodes: seasonTotalEpisodes,
		})
	})
	cc.OnError(func(r *colly.Response, err error) {
		log.Println("Error collecting searson info list.", err.Error())
	})

	tvDetailURL := fmt.Sprintf("%s/tv/%s/seasons", t.BaseURL, mediaID)
	log.Println("Visiting season info url:", tvDetailURL)

	cc.Visit(tvDetailURL)

	cc.Wait()

	return out
}

func (t TmdbScrapper) SearchMovie(in models.ClearFileEntry) ([]models.MovieResult, error) {
	searchString := t.GetSearchableString(in)
	out := []models.MovieResult{}

	log.Printf("Searching for movie: %s\n", searchString)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Accept-Language", "en-US")
	})

	c.OnHTML(".search_results.movie div.card.v4", func(h *colly.HTMLElement) {
		mediaInfo := t.ScrapMediaInfoListFromCollyElement(h)
		log.Println("MediaInfo:", mediaInfo)

		result := models.MovieResult{
			MediaInfo: mediaInfo,
		}
		out = append(out, result)
	})

	// prepare for visiting the url
	pathEscape := url.PathEscape("query=" + searchString)
	visitURL := fmt.Sprintf("%s/search/movie?%s", t.BaseURL, pathEscape)

	log.Println("Visiting URL:", visitURL)
	err := c.Visit(visitURL)
	c.Wait()

	if err != nil {
		return nil, errors.New("Cannot visit search movie url. " + err.Error())
	}

	log.Printf("Total results found for movie %s: %d\n", in.Name, len(out))
	return out, nil
}

func (t TmdbScrapper) GetProperMediaName(mediaName string, year int, mediaID string) string {
	out := mediaName
	if year > 0 {
		out += fmt.Sprintf(" (%d)", year)
	}
	if len(mediaID) > 0 {
		out += fmt.Sprintf(" [tmdbid-%s]", mediaID)
	}
	return out
}

func (t TmdbScrapper) GetProperEpisodeName(mediaName string, season, episode int) string {
	// returns "Season XX/mediaName SXXEYY"
	if season == -1 || episode == -1 {
		return ""
	}
	return fmt.Sprintf("Season %02d/%s S%02dE%02d", season, mediaName, season, episode)
}
