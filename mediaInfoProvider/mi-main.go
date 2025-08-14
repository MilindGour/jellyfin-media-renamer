package mediainfoprovider

type MediaInfoProvider interface {
	SearchMovies(term string) []MovieResult
	SearchTVShows(term string) []TVResult
}
