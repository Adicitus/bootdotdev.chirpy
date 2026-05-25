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
    $1
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ClearUsers :exec

DELETE FROM users;