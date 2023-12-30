package handballnet

import (
	"fmt"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapSeason(round round) (*model.Season, error) {
	var season model.Season

	if err := db.Get().Where(&model.Season{UID: round.ID}).Attrs(model.Season{
		StartDate: mapTime(round.StartsAt),
		EndDate:   mapTime(round.EndsAt),
	}).FirstOrCreate(&season).Error; err != nil {
		return nil, fmt.Errorf("error creating season: %s", err)
	}

	return &season, nil
}
