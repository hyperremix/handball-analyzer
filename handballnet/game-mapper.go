package handballnet

import (
	"fmt"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapGame(game game, league model.League, homeTeam model.Team, awayTeam model.Team, referees []model.Referee) (model.Game, error) {
	var mappedGame model.Game

	if err := db.Get().Where(&model.Game{UID: game.ID}).Attrs(model.Game{
		Date:       mapTime(game.StartsAt),
		LeagueID:   league.ID,
		HomeTeamID: homeTeam.ID,
		AwayTeamID: awayTeam.ID,
		Referees:   referees,
	}).FirstOrCreate(&mappedGame).Error; err != nil {
		return model.Game{}, fmt.Errorf("error creating referee: %s", err)
	}

	return mappedGame, nil
}
