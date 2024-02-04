package handballnet

import (
	"github.com/hyperremix/handball-analyzer/db"
)

func mapLeague(tournamentData tournament, seasonID int64) db.UpsertLeagueParams {
	return db.UpsertLeagueParams{
		Uid:      tournamentData.ID,
		SeasonID: seasonID,
		Name:     tournamentData.Name,
	}
}
