-- name: GetLeague :one
SELECT * FROM leagues WHERE id = $1 LIMIT 1;

-- name: ListSeasonLeagues :many
SELECT * FROM leagues WHERE season_id = $1 ORDER BY name;

-- name: UpsertLeague :one
INSERT INTO
    leagues (uid, season_id, name)
VALUES ($1, $2, $3) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
