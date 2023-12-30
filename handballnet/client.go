package handballnet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
	"github.com/rs/zerolog/log"
)

var baseUrl = "https://www.handball.net/_next/data/0h-TeKlhev7_XByVoOBBT/leagues"

func getSchedule(leagueId string) (*scheduleResponse, error) {
	url := fmt.Sprintf("%s/%s/schedule.json", baseUrl, leagueId)
	log.Info().Msgf("Fetching schedule from %s", url)

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	scheduleResponse := scheduleResponse{}
	jsonErr := json.Unmarshal(body, &scheduleResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &scheduleResponse, nil
}

func getNewSeasonData(leagueId string, roundId string, gameId string) (*model.Season, error) {
	url := fmt.Sprintf("%s/%s/schedule/rounds/%s/games/%s.json", baseUrl, leagueId, roundId, gameId)
	log.Info().Msgf("Fetching game from %s", url)

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		log.Error().Err(reqErr).Msgf("Could not create get game request for id=%s", gameId)
		return nil, reqErr
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		log.Error().Err(resErr).Msgf("Could not get game request for id=%s", gameId)
		return nil, resErr
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Error().Err(readErr).Msgf("Could not read get game response body for id=%s", gameId)
		return nil, readErr
	}

	gameResponse := gameResponse{}
	jsonErr := json.Unmarshal(body, &gameResponse)
	if jsonErr != nil {
		log.Error().Err(jsonErr).Msgf("Could not unmarshal get game response body for id=%s", gameId)
		return nil, jsonErr
	}

	season, err := mapSeason(gameResponse.PageProps.Game.Round)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map season for id=%s", gameId)
		return nil, err
	}

	league, err := mapLeague(gameResponse.PageProps.Game.Tournament, season.ID)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map league for id=%s", gameId)
		return nil, err
	}

	homeTeam, err := mapTeam(gameResponse.PageProps.Game.HomeTeam, league.ID)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map home team for id=%s", gameId)
		return nil, err
	}

	homeTeamMembers, err := mapTeamMembers(gameResponse.PageProps.Lineup.Home, gameResponse.PageProps.Lineup.HomeOfficials, homeTeam.ID)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map home team members for id=%s", gameId)
		return nil, err
	}

	awayTeam, err := mapTeam(gameResponse.PageProps.Game.AwayTeam, homeTeam.LeagueID)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map away team for id=%s", gameId)
		return nil, err
	}

	awayTeamMembers, err := mapTeamMembers(gameResponse.PageProps.Lineup.Away, gameResponse.PageProps.Lineup.AwayOfficials, awayTeam.ID)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map away team members for id=%s", gameId)
		return nil, err
	}

	referees, err := mapReferees(gameResponse.PageProps.Lineup.Referees)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map referees for id=%s", gameId)
		return nil, err
	}

	homeTeam.Members = homeTeamMembers
	awayTeam.Members = awayTeamMembers

	game, err := mapGame(gameResponse.PageProps.Game, league, homeTeam, awayTeam, referees)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map game for id=%s", gameId)
		return nil, err
	}

	gameScores, err := mapManyGameEvents(game, homeTeam, awayTeam, gameResponse.PageProps.Events)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map game events for id=%s", gameId)
		return nil, err
	}

	game.HalftimeScoreHome = gameScores.HalftimeHome
	game.HalftimeScoreAway = gameScores.HalftimeAway
	game.FulltimeScoreHome = gameScores.FulltimeHome
	game.FulltimeScoreAway = gameScores.FulltimeAway

	db.Get().Save(&game)

	return season, nil
}
