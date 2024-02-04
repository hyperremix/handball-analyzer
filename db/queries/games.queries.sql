-- name: GetGame :one
SELECT * FROM games WHERE id = $1 LIMIT 1;

-- name: ListLeagueGames :many
SELECT * FROM games WHERE league_id = $1 ORDER BY date DESC;

-- name: ListGamesByLeagueUid :many
SELECT games.*
FROM games
INNER JOIN leagues ON games.league_id = leagues.id
WHERE leagues.uid = $1
ORDER BY games.date DESC;

-- name: UpsertGame :one
INSERT INTO
    games (uid, league_id, date, home_team_id, away_team_id)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uid)
DO
UPDATE
SET
    updated_at = NOW() RETURNING *;

-- name: UpdateGameScores :exec
UPDATE
    games
SET
    halftime_home_score = $2,
    halftime_away_score = $3,
    fulltime_home_score = $4,
    fulltime_away_score = $5,
    updated_at = NOW()
WHERE
    id = $1;