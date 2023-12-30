package handballnet

import (
	"strconv"
	"strings"
	"time"
)

func mapTime(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

func mapElapsedSeconds(minutesAndSecondsString string) (uint, error) {
	minutesAndSeconds := strings.Split(minutesAndSecondsString, ":")

	minutes, err := strconv.Atoi(minutesAndSeconds[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(minutesAndSeconds[1])
	if err != nil {
		return 0, err
	}

	return uint(minutes*60 + seconds), nil
}
