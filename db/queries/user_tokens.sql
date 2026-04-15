-- name: CreateUserToken :one
INSERT INTO user_tokens (user_id, access_token, refresh_token, access_token_expires_at, refresh_token_expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserTokenByAccessToken :one
SELECT * FROM user_tokens
WHERE access_token = $1 AND is_active = TRUE;

-- name: GetUserTokenByRefreshToken :one
SELECT * FROM user_tokens
WHERE refresh_token = $1 AND is_active = TRUE;

-- name: DeactivateToken :exec
UPDATE user_tokens SET is_active = FALSE WHERE access_token = $1;

-- name: DeactivateAllUserTokens :exec
UPDATE user_tokens SET is_active = FALSE WHERE user_id = $1;
