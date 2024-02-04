-- name: ListGameRedCards :many
SELECT * FROM game_event_red_cards WHERE game_id = $1 ORDER BY elapsed_seconds;

-- name: UpsertGameEventRedCard :one
INSERT INTO
    game_event_red_cards (uid, game_id, team_id, team_member_id, daytime, elapsed_seconds)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
