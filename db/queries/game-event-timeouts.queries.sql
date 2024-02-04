-- name: ListGameTimeouts :many
SELECT * FROM game_event_timeouts WHERE game_id = $1 ORDER BY elapsed_seconds;

-- name: UpsertGameEventTimeout :one
INSERT INTO
    game_event_timeouts (uid, game_id, team_id, daytime, elapsed_seconds)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
