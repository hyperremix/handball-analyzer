package handballnet

import (
	"fmt"
	"strings"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
	"gorm.io/gorm/clause"
)

type gameEvents struct {
	blueCards   []model.GameEventBlueCard
	goals       []model.GameEventGoal
	penalties   []model.GameEventPenalty
	redCards    []model.GameEventRedCard
	sevenMeters []model.GameEventSevenMeters
	timeouts    []model.GameEventTimeout
	yellowCards []model.GameEventYellowCard
}

type gameScores struct {
	HalftimeHome uint
	HalftimeAway uint
	FulltimeHome uint
	FulltimeAway uint
}

func mapManyGameEvents(game model.Game, homeTeam model.Team, awayTeam model.Team, events []event) (gameScores, error) {
	allGameEvents := gameEvents{
		blueCards:   []model.GameEventBlueCard{},
		goals:       []model.GameEventGoal{},
		penalties:   []model.GameEventPenalty{},
		redCards:    []model.GameEventRedCard{},
		sevenMeters: []model.GameEventSevenMeters{},
		timeouts:    []model.GameEventTimeout{},
		yellowCards: []model.GameEventYellowCard{},
	}
	var scores gameScores

	for _, event := range events {
		if event.Type == eventTypeStopPeriod {
			scoreHome, scoreAway, isHalftime, err := mapStopPeriod(event)
			if err != nil {
				return gameScores{}, err
			}

			if isHalftime {
				scores.HalftimeHome = scoreHome
				scores.HalftimeAway = scoreAway
				continue
			}

			scores.FulltimeHome = scoreHome
			scores.FulltimeAway = scoreAway
			continue
		}

		team, err := getTeam(homeTeam, awayTeam, event)
		if err != nil {
			return gameScores{}, err
		}

		if event.Type == eventTypeTimeout {
			timeout, err := mapGameEventTimeout(game, team, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.timeouts = append(allGameEvents.timeouts, timeout)
			continue
		}

		teamMember, err := getTeamMember(team.Members, event)
		if err != nil && event.Type == eventTypeWarning {
			yellowCard, err := mapGameEventYellowCard(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.yellowCards = append(allGameEvents.yellowCards, yellowCard)
			continue
		}

		if err != nil && event.Type == eventTypeTwoMinutePenalty {
			penalty, err := mapGameEventPenalty(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.penalties = append(allGameEvents.penalties, penalty)
			continue
		}

		if err != nil {
			return gameScores{}, err
		}

		switch event.Type {
		case eventTypeGoal:
			goal, err := mapGameEventGoal(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.goals = append(allGameEvents.goals, goal)
			continue
		case eventTypeSevenMeterGoal:
			sevenMeterGoal, err := mapGameEventSevenMeters(game, team, teamMember, event, true)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.sevenMeters = append(allGameEvents.sevenMeters, sevenMeterGoal)
			continue
		case eventTypeTwoMinutePenalty:
			penalty, err := mapGameEventPenalty(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.penalties = append(allGameEvents.penalties, penalty)
			continue
		case eventTypeSevenMeterMissed:
			sevenMeterMissed, err := mapGameEventSevenMeters(game, team, teamMember, event, false)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.sevenMeters = append(allGameEvents.sevenMeters, sevenMeterMissed)
			continue
		case eventTypeWarning:
			yellowCard, err := mapGameEventYellowCard(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.yellowCards = append(allGameEvents.yellowCards, yellowCard)
			continue
		case eventTypeDisqualification:
			redCard, err := mapGameEventRedCard(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.redCards = append(allGameEvents.redCards, redCard)
			continue
		case eventTypeDisqualificationWithBlueCard:
			blueCard, err := mapGameEventBlueCard(game, team, teamMember, event)
			if err != nil {
				return gameScores{}, err
			}

			allGameEvents.blueCards = append(allGameEvents.blueCards, blueCard)
			continue
		default:
			return gameScores{}, fmt.Errorf(fmt.Sprintf("unknown event type %v", event.Type))
		}
	}

	if len(allGameEvents.blueCards) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.blueCards).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating blue cards: %s", err)
		}
	}

	if len(allGameEvents.goals) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.goals).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating goals: %s", err)
		}
	}

	if len(allGameEvents.penalties) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.penalties).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating penalties: %s", err)
		}
	}

	if len(allGameEvents.redCards) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.redCards).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating red cards: %s", err)
		}
	}

	if len(allGameEvents.sevenMeters) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.sevenMeters).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating seven meters: %s", err)
		}
	}

	if len(allGameEvents.timeouts) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.timeouts).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating timeouts: %s", err)
		}
	}

	if len(allGameEvents.yellowCards) != 0 {
		if err := db.Get().Clauses(clause.OnConflict{DoNothing: true}).Create(&allGameEvents.yellowCards).Error; err != nil {
			return gameScores{}, fmt.Errorf("error creating yellow cards: %s", err)
		}
	}

	return scores, nil
}

func mapStopPeriod(event event) (uint, uint, bool, error) {
	if event.Score == "" {
		return 0, 0, false, fmt.Errorf("could not map score for event %v", event)
	}

	scoreHome, scoreAway, err := mapScore(event.Score)
	if err != nil {
		return 0, 0, false, fmt.Errorf("could not map score for event %v", event)
	}

	isHalftime := strings.Contains(event.Message, "Spielstand 1. Halbzeit")

	return scoreHome, scoreAway, isHalftime, nil
}

func getTeam(homeTeam model.Team, awayTeam model.Team, event event) (model.Team, error) {
	if event.Team == "Home" {
		return homeTeam, nil
	}

	if event.Team == "Away" {
		return awayTeam, nil
	}

	return model.Team{}, fmt.Errorf("unknown team %s", event.Team)
}

func getTeamMember(teamMembers []model.TeamMember, event event) (model.TeamMember, error) {
	for _, teamMember := range teamMembers {
		if strings.Contains(event.Message, teamMember.Name) {
			return teamMember, nil
		}

		if strings.Contains(event.Message, fmt.Sprintf("%v.", teamMember.Number)) {
			return teamMember, nil
		}
	}

	return model.TeamMember{}, fmt.Errorf("could not find team member within message \"%s\" %v", event.Message, teamMembers)
}

func mapGameEventGoal(game model.Game, team model.Team, teamMember model.TeamMember, event event) (model.GameEventGoal, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypeGoal, event)
	if err != nil {
		return model.GameEventGoal{}, err
	}

	scoreHome, scoreAway, err := mapScore(event.Score)
	if err != nil {
		return model.GameEventGoal{}, fmt.Errorf("could not map score for event %v", event)
	}

	return model.GameEventGoal{
		BaseGameEvent: baseGameEvent,
		TeamMemberID:  teamMember.ID,
		ScoreHome:     scoreHome,
		ScoreAway:     scoreAway,
	}, nil
}

func mapGameEventPenalty(game model.Game, team model.Team, teamMember model.TeamMember, event event) (model.GameEventPenalty, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypePenalty, event)
	if err != nil {
		return model.GameEventPenalty{}, err
	}

	return model.GameEventPenalty{
		BaseGameEvent: baseGameEvent,
		TeamMemberID:  teamMember.ID,
	}, nil
}

func mapGameEventTimeout(game model.Game, team model.Team, event event) (model.GameEventTimeout, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypeTimeout, event)
	if err != nil {
		return model.GameEventTimeout{}, err
	}

	return model.GameEventTimeout{
		BaseGameEvent: baseGameEvent,
	}, nil
}

func mapGameEventSevenMeters(game model.Game, team model.Team, teamMember model.TeamMember, event event, isGoal bool) (model.GameEventSevenMeters, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypeSevenMeters, event)
	if err != nil {
		return model.GameEventSevenMeters{}, err
	}

	scoreHome, scoreAway, err := mapScore(event.Score)
	if err != nil {
		return model.GameEventSevenMeters{}, fmt.Errorf("could not map score for event %v", event)
	}

	return model.GameEventSevenMeters{
		BaseGameEvent: baseGameEvent,
		TeamMemberID:  teamMember.ID,
		IsGoal:        isGoal,
		ScoreHome:     scoreHome,
		ScoreAway:     scoreAway,
	}, nil
}

func mapGameEventYellowCard(game model.Game, team model.Team, teamMember model.TeamMember, event event) (model.GameEventYellowCard, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypeYellowCard, event)
	if err != nil {
		return model.GameEventYellowCard{}, err
	}

	return model.GameEventYellowCard{
		BaseGameEvent: baseGameEvent,
		TeamMemberID:  teamMember.ID,
	}, nil
}

func mapGameEventRedCard(game model.Game, team model.Team, teamMember model.TeamMember, event event) (model.GameEventRedCard, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypeRedCard, event)
	if err != nil {
		return model.GameEventRedCard{}, err
	}

	return model.GameEventRedCard{
		BaseGameEvent: baseGameEvent,
		TeamMemberID:  teamMember.ID,
	}, nil
}

func mapGameEventBlueCard(game model.Game, team model.Team, teamMember model.TeamMember, event event) (model.GameEventBlueCard, error) {
	baseGameEvent, err := mapBaseGameEvent(game, team, model.GameEventTypeBlueCard, event)
	if err != nil {
		return model.GameEventBlueCard{}, err
	}

	return model.GameEventBlueCard{
		BaseGameEvent: baseGameEvent,
		TeamMemberID:  teamMember.ID,
	}, nil
}

func mapBaseGameEvent(game model.Game, team model.Team, gameEventType model.GameEventType, event event) (model.BaseGameEvent, error) {
	elapsedSeconds, err := mapElapsedSeconds(event.Time)
	if err != nil {
		return model.BaseGameEvent{}, fmt.Errorf("could not map elapsed seconds for event %v", event)
	}

	return model.BaseGameEvent{
		UID:            fmt.Sprintf("%s.%v", game.UID, event.ID),
		GameID:         game.ID,
		Type:           gameEventType,
		Daytime:        mapTime(event.Timestamp),
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
	}, nil
}
