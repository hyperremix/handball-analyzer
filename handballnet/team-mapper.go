package handballnet

import (
	"fmt"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapTeam(team team, leagueID uint) (model.Team, error) {
	var mappedTeam model.Team

	if err := db.Get().Where(&model.Team{UID: team.ID}).Attrs(model.Team{
		LeagueID: leagueID,
		Name:     team.Name,
	}).FirstOrCreate(&mappedTeam).Error; err != nil {
		return model.Team{}, fmt.Errorf("error creating team: %s", err)
	}

	return mappedTeam, nil
}
