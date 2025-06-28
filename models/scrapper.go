package models

type MovieResult struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	YearOfRelease int    `json:"yearOfRelease"`
	ThumnailUrl   string `json:"thumnailUrl"`
}

type TVResult struct {
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	YearOfRelease int          `json:"yearOfRelease"`
	ThumnailUrl   string       `json:"thumnailUrl"`
	TotalSeasons  int          `json:"totalSeasons"`
	Seasons       []SeasonInfo `json:"seasons"`
}

type SeasonInfo struct {
	Number        int `json:"number"`
	TotalEpisodes int `json:"totalEpisodes"`
}
