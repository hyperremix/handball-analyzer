package services

import (
	"fmt"
	"sort"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func MapToLeagueListResponse(seasonRow db.Season, leagueRows []db.League) model.LeagueListResponse {
	var response model.LeagueListResponse

	response.SeasonID = seasonRow.ID
	response.SeasonName = mapSeasonName(seasonRow)

	for _, league := range leagueRows {
		response.Leagues = append(response.Leagues, model.LeagueResponse{
			ID:   league.ID,
			Name: league.Name,
		})
	}

	return response
}

func mapSeasonName(seasonRow db.Season) string {
	return fmt.Sprintf("Season %s/%s", seasonRow.StartDate.Time.Format("2006"), seasonRow.EndDate.Time.Format("2006"))
}

func MapToLeagueDetailsResponse(leagueRow db.League, gameRows []db.Game, teamRows []db.Team) model.LeagueDetailsResponse {
	teamStatsMap := getTeamStatsMap(teamRows)

	var details model.LeagueDetailsResponse

	details.Name = leagueRow.Name

	for _, game := range gameRows {
		if entry, ok := teamStatsMap[game.HomeTeamID]; ok {
			if game.FulltimeHomeScore > game.FulltimeAwayScore {
				entry.GamesWon++
				entry.Points += 2
			} else if game.FulltimeHomeScore == game.FulltimeAwayScore {
				entry.GamesDrawn++
				entry.Points++
			} else {
				entry.GamesLost++
				entry.PointsAgainst += 2
			}

			entry.GamesPlayed++
			entry.Goals += int(game.FulltimeHomeScore)
			entry.GoalsAgainst += int(game.FulltimeAwayScore)
			entry.GoalDifference += int(game.FulltimeHomeScore) - int(game.FulltimeAwayScore)

			teamStatsMap[game.HomeTeamID] = entry
		}

		if entry, ok := teamStatsMap[game.AwayTeamID]; ok {
			if game.FulltimeAwayScore > game.FulltimeHomeScore {
				entry.GamesWon++
				entry.Points += 2
			} else if game.FulltimeAwayScore == game.FulltimeHomeScore {
				entry.GamesDrawn++
				entry.Points++
			} else {
				entry.GamesLost++
				entry.PointsAgainst += 2
			}

			entry.GamesPlayed++
			entry.Goals += int(game.FulltimeAwayScore)
			entry.GoalsAgainst += int(game.FulltimeHomeScore)
			entry.GoalDifference += int(game.FulltimeAwayScore) - int(game.FulltimeHomeScore)

			teamStatsMap[game.AwayTeamID] = entry
		}
	}

	details.TeamStats = statsMapToSortedArray(teamStatsMap)

	return details
}

func getTeamStatsMap(teams []db.Team) map[int64]model.TeamStatsResponse {
	teamStatsMap := make(map[int64]model.TeamStatsResponse)

	for _, team := range teams {
		teamStatsMap[team.ID] = model.TeamStatsResponse{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	return teamStatsMap
}

func statsMapToSortedArray(statsMap map[int64]model.TeamStatsResponse) []model.TeamStatsResponse {
	var statsArray []model.TeamStatsResponse

	for _, stats := range statsMap {
		statsArray = append(statsArray, stats)
	}

	sort.Slice(statsArray, func(i, j int) bool {
		if statsArray[i].Points != statsArray[j].Points {
			return statsArray[i].Points > statsArray[j].Points
		}

		if statsArray[i].PointsAgainst != statsArray[j].PointsAgainst {
			return statsArray[i].PointsAgainst < statsArray[j].PointsAgainst
		}

		return statsArray[i].GoalDifference > statsArray[j].GoalDifference
	})

	return statsArray
}
