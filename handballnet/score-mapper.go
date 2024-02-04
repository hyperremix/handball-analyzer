package handballnet

import (
	"strconv"
	"strings"
)

func mapScore(scoreString string) (int32, int32, error) {
	homeAndAway := strings.Split(scoreString, ":")

	homeScore, err := strconv.Atoi(homeAndAway[0])
	if err != nil {
		return 0, 0, err
	}

	awayScore, err := strconv.Atoi(homeAndAway[1])
	if err != nil {
		return 0, 0, err
	}

	return int32(homeScore), int32(awayScore), nil
}
