package model

import "time"

type SeasonResponse struct {
	ID        int64
	StartDate time.Time
	EndDate   time.Time
	Name      string
}
