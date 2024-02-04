-- name: ListGameGoals :many
SELECT * FROM game_event_goals WHERE game_id = $1 ORDER BY elapsed_seconds;

-- name: UpsertGameEventGoal :one
INSERT INTO
    game_event_goals (uid, game_id, team_id, team_member_id, daytime, elapsed_seconds, score_home, score_away)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
