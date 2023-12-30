package handballnet

import (
	"slices"
	"time"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

var leagueUID = "handball4all.hamburg.m-ll-122_hhv"

func ProcessNewGames() error {
	scheduleResponse, err := getSchedule(leagueUID)
	if err != nil {
		return err
	}

	var league model.League
	var existingGames []model.Game

	if err := db.Get().Model(&league).Where(&model.League{UID: leagueUID}).Association("Games").Find(&existingGames); err != nil {
		return err
	}

	for _, game := range scheduleResponse.PageProps.Schedule.Data {
		startsAt := mapTime(game.StartsAt)

		index := slices.IndexFunc(existingGames, func(g model.Game) bool {
			return g.UID == game.ID
		})

		if index != -1 || startsAt.After(time.Now().Add(-2*time.Hour)) {
			continue
		}

		season, err := getNewSeasonData(leagueUID, game.Round.ID, game.ID)
		if err != nil {
			return err
		}

		if err := db.Get().Save(&season).Error; err != nil {
			return err
		}
	}

	return nil
}
