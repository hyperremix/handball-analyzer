package handballnet

import (
	"strconv"
	"strings"
)

func mapScore(scoreString string) (uint, uint, error) {
	homeAndAway := strings.Split(scoreString, ":")

	homeScore, err := strconv.Atoi(homeAndAway[0])
	if err != nil {
		return 0, 0, err
	}

	awayScore, err := strconv.Atoi(homeAndAway[1])
	if err != nil {
		return 0, 0, err
	}

	return uint(homeScore), uint(awayScore), nil
}
