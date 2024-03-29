// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: game-event-timeouts.queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listGameTimeouts = `-- name: ListGameTimeouts :many
SELECT id, uid, game_id, team_id, daytime, elapsed_seconds, created_at, updated_at, deleted_at FROM game_event_timeouts WHERE game_id = $1 ORDER BY elapsed_seconds
`

func (q *Queries) ListGameTimeouts(ctx context.Context, gameID int64) ([]GameEventTimeout, error) {
	rows, err := q.db.Query(ctx, listGameTimeouts, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GameEventTimeout
	for rows.Next() {
		var i GameEventTimeout
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.GameID,
			&i.TeamID,
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

const upsertGameEventTimeout = `-- name: UpsertGameEventTimeout :one
INSERT INTO
    game_event_timeouts (uid, game_id, team_id, daytime, elapsed_seconds)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING id, uid, game_id, team_id, daytime, elapsed_seconds, created_at, updated_at, deleted_at
`

type UpsertGameEventTimeoutParams struct {
	Uid            string
	GameID         int64
	TeamID         int64
	Daytime        pgtype.Timestamptz
	ElapsedSeconds int32
}

func (q *Queries) UpsertGameEventTimeout(ctx context.Context, arg UpsertGameEventTimeoutParams) (GameEventTimeout, error) {
	row := q.db.QueryRow(ctx, upsertGameEventTimeout,
		arg.Uid,
		arg.GameID,
		arg.TeamID,
		arg.Daytime,
		arg.ElapsedSeconds,
	)
	var i GameEventTimeout
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.GameID,
		&i.TeamID,
		&i.Daytime,
		&i.ElapsedSeconds,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
