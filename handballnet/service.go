package handballnet

import (
	"context"
	"time"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/environment"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var leagueUID = "handball4all.hamburg.m-ll-122_hhv"

func ProcessNewGames() error {
	scheduleResponse, err := getSchedule(leagueUID)
	if err != nil {
		return err
	}

	// var existingGames []db.Game

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, environment.DB_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// existingGames, err = queries.ListGamesByLeagueUid(ctx, leagueUID)
	// if err != nil {
	// 	return err
	// }

	for _, game := range scheduleResponse.PageProps.Schedule.Data {
		startsAt := mapTime(game.StartsAt)

		// index := slices.IndexFunc(existingGames, func(g db.Game) bool {
		// 	return g.Uid == game.ID
		// })

		if startsAt.After(time.Now().Add(-2 * time.Hour)) {
			continue
		}

		gameResponse, err := getGame(leagueUID, game.Round.ID, game.ID)
		if err != nil {
			return err
		}

		err = processNewGame(ctx, queries, gameResponse)
		if err != nil {
			return err
		}
	}

	return nil
}

func processNewGame(ctx context.Context, queries *db.Queries, gameResponse *gameResponse) error {
	gameId := gameResponse.PageProps.Game.ID

	season, err := queries.UpsertSeason(ctx, mapSeason(gameResponse.PageProps.Game.Round))
	if err != nil {
		log.Error().Err(err).Msgf("Could not map season for id=%s", gameId)
		return err
	}

	league, err := queries.UpsertLeague(ctx, mapLeague(gameResponse.PageProps.Game.Tournament, season.ID))
	if err != nil {
		log.Error().Err(err).Msgf("Could not map league for id=%s", gameId)
		return err
	}

	homeTeam, err := queries.UpsertTeam(ctx, mapTeam(gameResponse.PageProps.Game.HomeTeam, league.ID))
	if err != nil {
		log.Error().Err(err).Msgf("Could not map home team for id=%s", gameId)
		return err
	}

	var homeTeamMembers []db.TeamMember
	upsertHomeTeamMembersParams := mapTeamMembers(gameResponse.PageProps.Lineup.Home, gameResponse.PageProps.Lineup.HomeOfficials, homeTeam.ID)
	for _, params := range upsertHomeTeamMembersParams {
		teamMember, err := queries.UpsertTeamMember(ctx, params)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map home team members for id=%s", gameId)
			return err
		}
		homeTeamMembers = append(homeTeamMembers, teamMember)
	}

	awayTeam, err := queries.UpsertTeam(ctx, mapTeam(gameResponse.PageProps.Game.AwayTeam, homeTeam.LeagueID))
	if err != nil {
		log.Error().Err(err).Msgf("Could not map away team for id=%s", gameId)
		return err
	}

	var awayTeamMembers []db.TeamMember
	upsertAwayTeamMembersParams := mapTeamMembers(gameResponse.PageProps.Lineup.Away, gameResponse.PageProps.Lineup.AwayOfficials, awayTeam.ID)
	for _, params := range upsertAwayTeamMembersParams {
		teamMember, err := queries.UpsertTeamMember(ctx, params)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map away team members for id=%s", gameId)
			return err
		}
		awayTeamMembers = append(awayTeamMembers, teamMember)
	}

	game, err := queries.UpsertGame(ctx, mapGame(gameResponse.PageProps.Game, league, homeTeam, awayTeam))
	if err != nil {
		log.Error().Err(err).Msgf("Could not map game for id=%s", gameId)
		return err
	}

	upsertRefereeParams := mapReferees(gameResponse.PageProps.Lineup.Referees)
	for _, params := range upsertRefereeParams {
		referee, err := queries.UpsertReferee(ctx, params)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map referees for id=%s", gameId)
			return err
		}

		err = queries.UpsertGameReferee(ctx, db.UpsertGameRefereeParams{
			GameID:    game.ID,
			RefereeID: referee.ID,
		})
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game referees for id=%s", gameId)
			return err
		}
	}

	gameScores, gameEvents, err := mapManyGameEvents(game, homeTeam, homeTeamMembers, awayTeam, awayTeamMembers, gameResponse.PageProps.Events)
	if err != nil {
		return err
	}

	err = queries.UpdateGameScores(ctx, gameScores)
	if err != nil {
		log.Error().Err(err).Msgf("Could not map game scores for id=%s", gameId)
		return err
	}

	err = db.UpsertGameEvents(ctx, queries, game.Uid, gameEvents.blueCards, gameEvents.goals, gameEvents.penalties, gameEvents.redCards, gameEvents.sevenMeters, gameEvents.timeouts, gameEvents.yellowCards)
	if err != nil {
		return err
	}

	return nil
}
