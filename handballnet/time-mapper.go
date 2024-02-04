package handballnet

import (
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func mapTime(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

func mapTimestamptz(timestamp int64) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  mapTime(timestamp),
		Valid: true,
	}
}

func mapElapsedSeconds(minutesAndSecondsString string) (int32, error) {
	minutesAndSeconds := strings.Split(minutesAndSecondsString, ":")

	minutes, err := strconv.Atoi(minutesAndSeconds[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(minutesAndSeconds[1])
	if err != nil {
		return 0, err
	}

	return int32(minutes*60 + seconds), nil
}
