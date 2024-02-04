package model

type LeagueListResponse struct {
	SeasonID   int64
	SeasonName string
	Leagues    []LeagueResponse
}

type LeagueResponse struct {
	ID   int64
	Name string
}
