package scrapper

type Scrapper interface {
	Scrap(url string, itemSel string, fieldMap map[string]string) (ScrapResultList, error)
}

type ScrapResultList []map[string]string
