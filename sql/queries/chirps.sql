-- name: CreateChirp :one
INSERT INTO chirps (
    id,
    created_at,
    updated_at,
    body,
    user_id
) VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    @body::TEXT,
    @user_id::UUID
) RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = @chirp_id::UUID;