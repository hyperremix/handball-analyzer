package model

type GameEventGoal struct {
	BaseGameEvent
	TeamMemberID uint
	ScoreHome    uint
	ScoreAway    uint
}
