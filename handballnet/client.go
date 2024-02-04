package handballnet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var baseUrl = "https://www.handball.net/_next/data/a5HpojOeVGw3jS67iJjiH/leagues"

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

func getGame(leagueId string, roundId string, gameId string) (*gameResponse, error) {
	url := fmt.Sprintf("%s/%s/schedule/rounds/%s/games/%s.json", baseUrl, leagueId, roundId, gameId)
	log.Info().Msgf("Fetching game from %s", url)

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msgf("Could not create get game request for id=%s", gameId)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get game request for id=%s", gameId)
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Could not read get game response body for id=%s", gameId)
		return nil, err
	}

	gameResponse := gameResponse{}
	err = json.Unmarshal(body, &gameResponse)
	if err != nil {
		log.Error().Err(err).Msgf("Could not unmarshal get game response body for id=%s", gameId)
		return nil, err
	}

	return &gameResponse, nil
}
