package model

type GameEventSevenMeters struct {
	BaseGameEvent
	TeamMemberID uint
	IsGoal       bool
	ScoreHome    uint
	ScoreAway    uint
}
