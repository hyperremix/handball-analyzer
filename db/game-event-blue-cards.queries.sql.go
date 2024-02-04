// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: game-event-blue-cards.queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listGameBlueCards = `-- name: ListGameBlueCards :many
SELECT id, uid, game_id, team_id, team_member_id, daytime, elapsed_seconds, created_at, updated_at, deleted_at FROM game_event_blue_cards WHERE game_id = $1 ORDER BY elapsed_seconds
`

func (q *Queries) ListGameBlueCards(ctx context.Context, gameID int64) ([]GameEventBlueCard, error) {
	rows, err := q.db.Query(ctx, listGameBlueCards, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GameEventBlueCard
	for rows.Next() {
		var i GameEventBlueCard
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.GameID,
			&i.TeamID,
			&i.TeamMemberID,
			&i.Daytime,
			&i.ElapsedSeconds,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertGameEventBlueCard = `-- name: UpsertGameEventBlueCard :one
INSERT INTO
    game_event_blue_cards (uid, game_id, team_id, team_member_id, daytime, elapsed_seconds)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING id, uid, game_id, team_id, team_member_id, daytime, elapsed_seconds, created_at, updated_at, deleted_at
`

type UpsertGameEventBlueCardParams struct {
	Uid            string
	GameID         int64
	TeamID         int64
	TeamMemberID   int64
	Daytime        pgtype.Timestamptz
	ElapsedSeconds int32
}

func (q *Queries) UpsertGameEventBlueCard(ctx context.Context, arg UpsertGameEventBlueCardParams) (GameEventBlueCard, error) {
	row := q.db.QueryRow(ctx, upsertGameEventBlueCard,
		arg.Uid,
		arg.GameID,
		arg.TeamID,
		arg.TeamMemberID,
		arg.Daytime,
		arg.ElapsedSeconds,
	)
	var i GameEventBlueCard
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.GameID,
		&i.TeamID,
		&i.TeamMemberID,
		&i.Daytime,
		&i.ElapsedSeconds,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
