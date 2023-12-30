package model

import (
	"gorm.io/gorm"
)

type TeamMember struct {
	gorm.Model
	UID    string `gorm:"uniqueIndex"`
	TeamID uint
	Name   string
	Number string
	Type   TeamMemberType
}
