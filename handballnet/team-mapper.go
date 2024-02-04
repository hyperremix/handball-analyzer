package handballnet

import (
	"github.com/hyperremix/handball-analyzer/db"
)

func mapTeam(team team, leagueID int64) db.UpsertTeamParams {
	return db.UpsertTeamParams{
		Uid:      team.ID,
		LeagueID: leagueID,
		Name:     team.Name,
	}
}
