package model

type GameEventType string

const (
	GameEventTypeGoal        GameEventType = "Goal"
	GameEventTypeSevenMeters GameEventType = "SevenMeters"
	GameEventTypePenalty     GameEventType = "Penalty"
	GameEventTypeTimeout     GameEventType = "Timeout"
	GameEventTypeYellowCard  GameEventType = "YellowCard"
	GameEventTypeRedCard     GameEventType = "RedCard"
	GameEventTypeBlueCard    GameEventType = "BlueCard"
)
