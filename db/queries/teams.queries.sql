-- name: ListLeagueTeams :many
SELECT * FROM teams WHERE league_id = $1 ORDER BY name;

-- name: UpsertTeam :one
INSERT INTO
    teams (uid, name, league_id)
VALUES ($1, $2, $3) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
