-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_email,
  refresh_token,
  is_revoked,
  created_at,
  expires_at
)
VALUES (
  $1,
  $2,
  $3,
  $4,
  NOW(),
  $5
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1
  AND is_revoked = FALSE;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;
