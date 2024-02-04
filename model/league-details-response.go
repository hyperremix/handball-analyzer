package model

type LeagueDetailsResponse struct {
	ID        int64
	Name      string
	TeamStats []TeamStatsResponse
}

type TeamStatsResponse struct {
	TeamID         int64
	TeamName       string
	GamesPlayed    int
	GamesWon       int
	GamesDrawn     int
	GamesLost      int
	Goals          int
	GoalsAgainst   int
	GoalDifference int
	Points         int
	PointsAgainst  int
}
