package model

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	UID       string `gorm:"uniqueIndex"`
	Leagues   []League
	StartDate time.Time
	EndDate   time.Time
}
