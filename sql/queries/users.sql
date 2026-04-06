-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (
  id,
  name,
  email,
  password_hash,
  role
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  name = $2,
  email = $3,
  password_hash = $4,
  role = $5,
  updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetDeletedUsers :many
SELECT * FROM users
WHERE deleted_at IS NOT NULL;
