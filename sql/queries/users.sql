-- name: CreateUser :one

INSERT INTO users (
    id,
    created_at,
    updated_at,
    email
) VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    @email::TEXT
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email;

-- name: RemoveUser :exec
DELETE FROM users WHERE id = @user_id;

-- name: ClearUsers :exec

DELETE FROM users;

-- name: SetEmail :one
UPDATE users SET updated_at = CURRENT_TIMESTAMP, email = @email::TEXT WHERE id = @user_id RETURNING *;