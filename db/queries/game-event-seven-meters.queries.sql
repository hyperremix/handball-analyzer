-- name: ListGameSevenMeters :many
SELECT * FROM game_event_seven_meters WHERE game_id = $1 ORDER BY elapsed_seconds;

-- name: UpsertGameEventSevenMeter :one
INSERT INTO
    game_event_seven_meters (uid, game_id, team_id, team_member_id, daytime, elapsed_seconds, is_goal, score_home, score_away)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
