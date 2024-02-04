-- name: UpsertTeamMember :one
INSERT INTO
    team_members (uid, team_id, name, number, type)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;
