-- name: RegisterUser :one
INSERT INTO users (email, name, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name       = COALESCE(sqlc.narg('name'), name),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserLastActive :one
UPDATE users
SET last_active_at = $2,
    updated_at     = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
