package model

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	UID      string `gorm:"uniqueIndex"`
	LeagueID uint
	Name     string
	Members  []TeamMember
}
