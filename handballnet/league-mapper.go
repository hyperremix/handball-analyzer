package handballnet

import (
	"fmt"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapLeague(tournamentData tournament, seasonID uint) (model.League, error) {
	var league model.League

	if err := db.Get().Where(&model.League{UID: tournamentData.ID}).Attrs(model.League{
		SeasonID: seasonID,
		Name:     tournamentData.Name,
	}).FirstOrCreate(&league).Error; err != nil {
		return model.League{}, fmt.Errorf("error creating league: %s", err)
	}

	return league, nil
}
