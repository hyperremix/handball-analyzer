package handballnet

import (
	"github.com/hyperremix/handball-analyzer/db"
)

func mapSeason(round round) db.UpsertSeasonParams {
	return db.UpsertSeasonParams{
		Uid:       round.ID,
		StartDate: mapTimestamptz(round.StartsAt),
		EndDate:   mapTimestamptz(round.EndsAt),
	}
}
