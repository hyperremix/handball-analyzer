package handballnet

import (
	"fmt"
	"strings"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type gameEvents struct {
	blueCards   []db.UpsertGameEventBlueCardParams
	goals       []db.UpsertGameEventGoalParams
	penalties   []db.UpsertGameEventPenaltyParams
	redCards    []db.UpsertGameEventRedCardParams
	sevenMeters []db.UpsertGameEventSevenMeterParams
	timeouts    []db.UpsertGameEventTimeoutParams
	yellowCards []db.UpsertGameEventYellowCardParams
}

func mapManyGameEvents(game db.Game, homeTeam db.Team, homeTeamMembers []db.TeamMember, awayTeam db.Team, awayTeamMembers []db.TeamMember, events []event) (db.UpdateGameScoresParams, gameEvents, error) {
	allGameEvents := gameEvents{
		blueCards:   []db.UpsertGameEventBlueCardParams{},
		goals:       []db.UpsertGameEventGoalParams{},
		penalties:   []db.UpsertGameEventPenaltyParams{},
		redCards:    []db.UpsertGameEventRedCardParams{},
		sevenMeters: []db.UpsertGameEventSevenMeterParams{},
		timeouts:    []db.UpsertGameEventTimeoutParams{},
		yellowCards: []db.UpsertGameEventYellowCardParams{},
	}
	updateGameScoresParams := db.UpdateGameScoresParams{
		ID: game.ID,
	}

	for _, event := range events {
		if event.Type == eventTypeStopPeriod {
			scoreHome, scoreAway, isHalftime, err := mapStopPeriod(event)
			if err != nil {
				return db.UpdateGameScoresParams{}, gameEvents{}, err
			}

			if isHalftime {
				updateGameScoresParams.HalftimeHomeScore = scoreHome
				updateGameScoresParams.HalftimeAwayScore = scoreAway
				continue
			}

			updateGameScoresParams.FulltimeHomeScore = scoreHome
			updateGameScoresParams.FulltimeAwayScore = scoreAway
			continue
		}

		team, teamMembers, err := getTeamAndMembers(homeTeam, homeTeamMembers, awayTeam, awayTeamMembers, event)
		if err != nil {
			return db.UpdateGameScoresParams{}, gameEvents{}, err
		}

		elapsedSeconds, err := mapElapsedSeconds(event.Time)
		if err != nil {
			return db.UpdateGameScoresParams{}, gameEvents{}, fmt.Errorf("could not map elapsed seconds for event %v", event)
		}

		gameEventUid := fmt.Sprintf("%s.%v", game.Uid, event.ID)
		daytime := mapTimestamptz(event.Timestamp)

		if event.Type == eventTypeTimeout {
			timeout := mapGameEventTimeout(game, team, event, gameEventUid, elapsedSeconds, daytime)
			allGameEvents.timeouts = append(allGameEvents.timeouts, timeout)
			continue
		}

		teamMember, err := getTeamMember(teamMembers, event)
		if err != nil && event.Type == eventTypeWarning {
			yellowCard := mapGameEventYellowCard(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)

			allGameEvents.yellowCards = append(allGameEvents.yellowCards, yellowCard)
			continue
		}

		if err != nil && event.Type == eventTypeTwoMinutePenalty {
			penalty := mapGameEventPenalty(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)

			allGameEvents.penalties = append(allGameEvents.penalties, penalty)
			continue
		}

		if err != nil {
			return db.UpdateGameScoresParams{}, gameEvents{}, err
		}

		switch event.Type {
		case eventTypeGoal:
			goal, err := mapGameEventGoal(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)
			if err != nil {
				return db.UpdateGameScoresParams{}, gameEvents{}, err
			}

			allGameEvents.goals = append(allGameEvents.goals, goal)
			continue
		case eventTypeSevenMeterGoal:
			sevenMeterGoal, err := mapGameEventSevenMeters(game, team, teamMember, event, true, gameEventUid, elapsedSeconds, daytime)
			if err != nil {
				return db.UpdateGameScoresParams{}, gameEvents{}, err
			}

			allGameEvents.sevenMeters = append(allGameEvents.sevenMeters, sevenMeterGoal)
			continue
		case eventTypeTwoMinutePenalty:
			penalty := mapGameEventPenalty(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)
			allGameEvents.penalties = append(allGameEvents.penalties, penalty)
			continue
		case eventTypeSevenMeterMissed:
			sevenMeterMissed, err := mapGameEventSevenMeters(game, team, teamMember, event, false, gameEventUid, elapsedSeconds, daytime)
			if err != nil {
				return db.UpdateGameScoresParams{}, gameEvents{}, err
			}

			allGameEvents.sevenMeters = append(allGameEvents.sevenMeters, sevenMeterMissed)
			continue
		case eventTypeWarning:
			yellowCard := mapGameEventYellowCard(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)
			allGameEvents.yellowCards = append(allGameEvents.yellowCards, yellowCard)
			continue
		case eventTypeDisqualification:
			redCard := mapGameEventRedCard(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)
			allGameEvents.redCards = append(allGameEvents.redCards, redCard)
			continue
		case eventTypeDisqualificationWithBlueCard:
			blueCard := mapGameEventBlueCard(game, team, teamMember, event, gameEventUid, elapsedSeconds, daytime)
			allGameEvents.blueCards = append(allGameEvents.blueCards, blueCard)
			continue
		default:
			return db.UpdateGameScoresParams{}, gameEvents{}, fmt.Errorf(fmt.Sprintf("unknown event type %v", event.Type))
		}
	}

	return updateGameScoresParams, allGameEvents, nil
}

func mapStopPeriod(event event) (int32, int32, bool, error) {
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

func getTeamAndMembers(homeTeam db.Team, homeTeamMembers []db.TeamMember, awayTeam db.Team, awayTeamMembers []db.TeamMember, event event) (db.Team, []db.TeamMember, error) {
	if event.Team == "Home" {
		return homeTeam, homeTeamMembers, nil
	}

	if event.Team == "Away" {
		return awayTeam, awayTeamMembers, nil
	}

	return db.Team{}, []db.TeamMember{}, fmt.Errorf("unknown team %s", event.Team)
}

func getTeamMember(teamMembers []db.TeamMember, event event) (*db.TeamMember, error) {
	for _, teamMember := range teamMembers {
		if strings.Contains(event.Message, teamMember.Name) {
			return &teamMember, nil
		}

		if strings.Contains(event.Message, fmt.Sprintf("%v.", teamMember.Number)) {
			return &teamMember, nil
		}
	}

	return &db.TeamMember{}, fmt.Errorf("could not find team member within message \"%s\" %v", event.Message, teamMembers)
}

func mapGameEventGoal(game db.Game, team db.Team, teamMember *db.TeamMember, event event, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) (db.UpsertGameEventGoalParams, error) {
	scoreHome, scoreAway, err := mapScore(event.Score)
	if err != nil {
		return db.UpsertGameEventGoalParams{}, fmt.Errorf("could not map score for event %v", event)
	}

	return db.UpsertGameEventGoalParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
		TeamMemberID:   teamMember.ID,
		ScoreHome:      scoreHome,
		ScoreAway:      scoreAway,
	}, nil
}

func mapGameEventPenalty(game db.Game, team db.Team, teamMember *db.TeamMember, event event, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) db.UpsertGameEventPenaltyParams {
	return db.UpsertGameEventPenaltyParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
		TeamMemberID:   mapInt8(teamMember.ID, teamMember == nil),
	}
}

func mapGameEventTimeout(game db.Game, team db.Team, event event, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) db.UpsertGameEventTimeoutParams {
	return db.UpsertGameEventTimeoutParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
	}
}

func mapGameEventSevenMeters(game db.Game, team db.Team, teamMember *db.TeamMember, event event, isGoal bool, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) (db.UpsertGameEventSevenMeterParams, error) {
	scoreHome, scoreAway, err := mapScore(event.Score)
	if err != nil {
		return db.UpsertGameEventSevenMeterParams{}, fmt.Errorf("could not map score for event %v", event)
	}

	return db.UpsertGameEventSevenMeterParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
		TeamMemberID:   teamMember.ID,
		IsGoal:         isGoal,
		ScoreHome:      scoreHome,
		ScoreAway:      scoreAway,
	}, nil
}

func mapGameEventYellowCard(game db.Game, team db.Team, teamMember *db.TeamMember, event event, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) db.UpsertGameEventYellowCardParams {

	return db.UpsertGameEventYellowCardParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
		TeamMemberID:   mapInt8(teamMember.ID, teamMember == nil),
	}
}

func mapInt8(i int64, valid bool) pgtype.Int8 {
	if !valid {
		return pgtype.Int8{Valid: false}
	}

	return pgtype.Int8{Int64: i, Valid: true}
}

func mapGameEventRedCard(game db.Game, team db.Team, teamMember *db.TeamMember, event event, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) db.UpsertGameEventRedCardParams {
	return db.UpsertGameEventRedCardParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
		TeamMemberID:   teamMember.ID,
	}
}

func mapGameEventBlueCard(game db.Game, team db.Team, teamMember *db.TeamMember, event event, gameEventUid string, elapsedSeconds int32, daytime pgtype.Timestamptz) db.UpsertGameEventBlueCardParams {
	return db.UpsertGameEventBlueCardParams{
		Uid:            gameEventUid,
		GameID:         game.ID,
		Daytime:        daytime,
		ElapsedSeconds: elapsedSeconds,
		TeamID:         team.ID,
		TeamMemberID:   teamMember.ID,
	}
}
