-- name: UpsertGameReferee :exec
INSERT INTO
    games_referees (game_id, referee_id)
VALUES ($1, $2) ON CONFLICT (game_id, referee_id)
DO
NOTHING;
