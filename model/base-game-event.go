package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseGameEvent struct {
	gorm.Model
	UID            string `gorm:"uniqueIndex"`
	Type           GameEventType
	GameID         uint
	Daytime        time.Time
	ElapsedSeconds uint
	TeamID         uint
}

type GameEvent interface {
	GetType() GameEventType
}

func (b *BaseGameEvent) GetType() GameEventType {
	return b.Type
}
