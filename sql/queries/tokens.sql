-- name: CreateToken :one
INSERT INTO tokens (
    user_id,
    created_at,
    updated_at,
    token,
    expires_at
) VALUES (
    @user_id::UUID,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    @token::TEXT,
    @expiration::TIMESTAMP
) RETURNING *;

-- name: GetToken :one
SELECT * FROM tokens WHERE token = @token;

-- name: RemoveToken :exec
DELETE FROM tokens WHERE token = @token;

-- name: RevokeToken :one
UPDATE tokens SET
    revoked_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE token = @token RETURNING *; 