package model

import (
	"gorm.io/gorm"
)

type Referee struct {
	gorm.Model
	UID  string `gorm:"uniqueIndex"`
	Name string
	Type RefereeType
}
