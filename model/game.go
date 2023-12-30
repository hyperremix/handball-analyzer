package model

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	UID               string `gorm:"uniqueIndex"`
	Date              time.Time
	LeagueID          uint
	HomeTeamID        uint
	AwayTeamID        uint
	HalftimeScoreHome uint
	HalftimeScoreAway uint
	FulltimeScoreHome uint
	FulltimeScoreAway uint
	Referees          []Referee `gorm:"many2many:games_referees;"`
}
