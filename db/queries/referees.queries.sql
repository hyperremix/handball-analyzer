-- name: UpsertReferee :one
INSERT INTO
    referees (uid, name, type)
VALUES ($1, $2, $3) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
