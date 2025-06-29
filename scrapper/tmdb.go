package scrapper

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gocolly/colly"
)

type TmdbScrapper struct {
	baseUrl string
}

func NewTmdbScrapper() *TmdbScrapper {
	out := TmdbScrapper{
		baseUrl: "https://www.themoviedb.org",
	}
	return &out
}

func (t TmdbScrapper) GetSearchableString(in models.ClearFileEntry) string {
	if in.Year == 0 {
		return in.Name
	}

	return fmt.Sprintf("%s y:%d", in.Name, in.Year)
}

func (t TmdbScrapper) SearchMovie(in models.ClearFileEntry) ([]models.MovieResult, error) {
	searchString := t.GetSearchableString(in)
	out := []models.MovieResult{}

	log.Printf("Seaching for movie: %s\n", searchString)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Accept-Language", "en-US")
	})

	c.OnHTML(".search_results.movie div.card.v4", func(h *colly.HTMLElement) {
		nameNodeString := h.ChildText(".title h2")
		nameTitleNodeString := h.ChildText(".title h2 .title")
		if len(nameTitleNodeString) > 0 {
			nameNodeString = strings.Replace(nameNodeString, nameTitleNodeString, "", 1)
		}
		nameNodeString = strings.Trim(nameNodeString, " ")

		result := models.MovieResult{
			Name:        nameNodeString,
			Description: h.ChildText(".overview"),
			ThumnailUrl: h.ChildAttr(".poster img", "src"),
		}
		releaseDateStr := h.ChildText(".release_date")
		releaseYear, err := util.ExtractYearFromString(releaseDateStr)
		if err != nil {
			releaseYear = 0
		}
		result.YearOfRelease = releaseYear
		out = append(out, result)
	})

	// prepare for visiting the url
	pathEscape := url.PathEscape("query=" + searchString)
	visitUrl := fmt.Sprintf("%s/search/movie?%s", t.baseUrl, pathEscape)

	err := c.Visit(visitUrl)
	c.Wait()

	if err != nil {
		return nil, errors.New("Cannot visit search movie url. " + err.Error())
	}

	return out, nil
}
