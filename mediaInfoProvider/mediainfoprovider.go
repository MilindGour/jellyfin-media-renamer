package mediainfoprovider

type MediaInfoProvider interface {
	SearchMovies(term string, year int) []MovieResult
	SearchTVShows(term string, year int) []TVResult
	GetProperDirectoryName(mediaInfo MediaInfo) string
	GetProperTVShowFilename(filename, showName string, season, episode int) string
}
