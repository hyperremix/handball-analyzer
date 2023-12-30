package model

import (
	"gorm.io/gorm"
)

type League struct {
	gorm.Model
	UID      string `gorm:"uniqueIndex"`
	SeasonID uint
	Name     string
	Games    []Game
	Teams    []Team
}
