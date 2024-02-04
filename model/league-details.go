package model

type LeagueDetails struct {
	Name      string
	TeamStats []TeamStats
}

type TeamStats struct {
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
