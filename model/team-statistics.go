package model

import (
	"gorm.io/gorm"
)

type TeamStatistics struct {
	gorm.Model
	TeamID           uint
	GameEvents       map[GameEventType]uint
	Wins             uint
	Losses           uint
	Draws            uint
	SevenMetersGoals uint
	ConcededGoals    uint
}
