package scrapper

import (
	"regexp"

	"github.com/MilindGour/jellyfin-media-renamer/network"
	"github.com/PuerkitoBio/goquery"
)

type GoQueryScrapper struct {
	htmlProvider network.HttpResponseProvider
}

func NewGoQueryScrapper() *GoQueryScrapper {
	return &GoQueryScrapper{
		htmlProvider: network.NewHttpResponse(),
	}
}
func NewMockGoQueryScrapper() *GoQueryScrapper {
	return &GoQueryScrapper{
		htmlProvider: network.NewMockResponse(),
	}
}

func (g *GoQueryScrapper) Scrap(url string, itemSel string, fieldMap map[string]string) (ScrapResultList, error) {
	res, err := g.htmlProvider.GetResponse(url)
	if err != nil {
		// Unable to get html
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		// Unable to parse html to doc
		return nil, err
	}

	out := ScrapResultList{}
	doc.Find(itemSel).Each(func(i int, s *goquery.Selection) {
		itemOutput := map[string]string{}
		for key, selector := range fieldMap {
			elSelector, attrSelector := g.splitAttribute(selector)
			val := ""
			el := s.Find(elSelector)
			if len(attrSelector) > 0 {
				attrVal, attrExist := el.Attr(attrSelector)
				if attrExist {
					val = attrVal
				}
			} else {
				val = el.Text()
			}
			itemOutput[key] = val
		}
		out = append(out, itemOutput)
	})

	return out, nil
}

// splitAttribute function splits attribute from selector if present.
// Otherwise it will return empty attribute name.
func (g *GoQueryScrapper) splitAttribute(selector string) (string, string) {
	outElement := selector
	outAttribute := ""

	lastCharacter := string(selector[len(selector)-1])
	if lastCharacter == "]" {
		// there is an attribute in last
		attrRe := regexp.MustCompile(`\[(\w+)\]$`)
		matchIndex := attrRe.FindStringIndex(selector)
		if matchIndex != nil {
			outElement = selector[0:matchIndex[0]]
			outAttribute = selector[matchIndex[0]+1 : matchIndex[1]-1]
		}
	}

	return outElement, outAttribute
}
