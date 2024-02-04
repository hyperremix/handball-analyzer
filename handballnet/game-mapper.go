package handballnet

import (
	"github.com/hyperremix/handball-analyzer/db"
)

func mapGame(game game, league db.League, homeTeam db.Team, awayTeam db.Team) db.UpsertGameParams {
	return db.UpsertGameParams{
		Uid:        game.ID,
		Date:       mapTimestamptz(game.StartsAt),
		LeagueID:   league.ID,
		HomeTeamID: homeTeam.ID,
		AwayTeamID: awayTeam.ID,
	}
}
