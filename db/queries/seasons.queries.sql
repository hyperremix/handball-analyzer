-- name: GetSeason :one
SELECT * FROM seasons WHERE id = $1 LIMIT 1;

-- name: ListSeasons :many
SELECT * FROM seasons ORDER BY start_date DESC;

-- name: UpsertSeason :one
INSERT INTO
    seasons (uid, start_date, end_date)
VALUES ($1, $2, $3) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
