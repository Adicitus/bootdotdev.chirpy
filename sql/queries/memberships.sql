-- name: ActivateMembership :one
UPDATE users SET is_chirpy_red = true, updated_at = CURRENT_TIMESTAMP WHERE id = @user_id RETURNING *;