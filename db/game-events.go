package db

import (
	"context"

	"github.com/rs/zerolog/log"
)

func UpsertGameEvents(
	ctx context.Context,
	queries *Queries,
	gameUid string,
	upsertGameEventBlueCardParams []UpsertGameEventBlueCardParams,
	upsertGameEventGoalsParams []UpsertGameEventGoalParams,
	upsertGameEventPenaltiesParams []UpsertGameEventPenaltyParams,
	upsertGameEventRedCardsParams []UpsertGameEventRedCardParams,
	upsertGameEventSevenMetersParams []UpsertGameEventSevenMeterParams,
	upsertGameEventTimeoutsParams []UpsertGameEventTimeoutParams,
	upsertGameEventYellowCardsParams []UpsertGameEventYellowCardParams,
) error {
	for _, gameEvent := range upsertGameEventBlueCardParams {
		_, err := queries.UpsertGameEventBlueCard(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event blue card for id=%s", gameUid)
			return err
		}
	}

	for _, gameEvent := range upsertGameEventGoalsParams {
		_, err := queries.UpsertGameEventGoal(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event goal for id=%s", gameUid)
			return err
		}
	}

	for _, gameEvent := range upsertGameEventPenaltiesParams {
		_, err := queries.UpsertGameEventPenalty(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event penalty for id=%s", gameUid)
			return err
		}
	}

	for _, gameEvent := range upsertGameEventRedCardsParams {
		_, err := queries.UpsertGameEventRedCard(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event red card for id=%s", gameUid)
			return err
		}
	}

	for _, gameEvent := range upsertGameEventSevenMetersParams {
		_, err := queries.UpsertGameEventSevenMeter(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event seven meter for id=%s", gameUid)
			return err
		}
	}

	for _, gameEvent := range upsertGameEventTimeoutsParams {
		_, err := queries.UpsertGameEventTimeout(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event timeout for id=%s", gameUid)
			return err
		}
	}

	for _, gameEvent := range upsertGameEventYellowCardsParams {
		_, err := queries.UpsertGameEventYellowCard(ctx, gameEvent)
		if err != nil {
			log.Error().Err(err).Msgf("Could not map game event yellow card for gameUid=%s, gameEventUid=%s, teamId=%v, teamMemberId=%v", gameUid, gameEvent.Uid, gameEvent.TeamID, gameEvent.TeamMemberID)
			return err
		}
	}

	return nil
}
