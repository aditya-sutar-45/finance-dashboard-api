-- name: CreateRecord :one
INSERT INTO records (
  id,
  user_id,
  amount,
  type,
  category,
  note,
  date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
