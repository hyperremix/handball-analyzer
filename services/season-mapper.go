package services

import (
	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func MapRowsToSeasonResponses(row []db.Season) []model.SeasonResponse {
	seasons := make([]model.SeasonResponse, len(row))
	for i := range row {
		seasons[i] = MapRowToSeasonResponse(row[i])
	}
	return seasons
}

func MapRowToSeasonResponse(row db.Season) model.SeasonResponse {
	return model.SeasonResponse{
		ID:        row.ID,
		StartDate: row.StartDate.Time,
		EndDate:   row.EndDate.Time,
		Name:      row.StartDate.Time.Format("2006") + "/" + row.EndDate.Time.Format("2006"),
	}
}
