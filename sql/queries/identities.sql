-- name: CreateIdentity :one
INSERT INTO identities (
    user_id,
    created_at,
    updated_at,
    auth
) VALUES (
    @user_id::UUID,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    @auth::TEXT
) RETURNING *;

-- name: GetIdentity :one
SELECT * FROM identities WHERE user_id = @user_id;

-- name: SetIdentity :one
UPDATE identities SET updated_at = CURRENT_TIMESTAMP, auth = @auth WHERE user_id = @user_id RETURNING *;