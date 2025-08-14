package jmrrenamer

type JmrRenamer interface {
	GetMediaNameAndYear(rawFilename string) MediaNameAndYear
	GetMediaSeasonAndEpisode(filePath string) MediaSeasonAndEpisode
}

type MediaNameAndYear struct {
	Name string
	Year int
}
type MediaSeasonAndEpisode struct {
	Season  int
	Episode int
}
