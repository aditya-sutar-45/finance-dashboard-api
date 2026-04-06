-- name: CreateRecord :one
INSERT INTO records (
  id,
  user_id,
  created_by,
  amount,
  type,
  category,
  note,
  date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListRecords :many
SELECT *
FROM records
WHERE user_id = sqlc.arg(user_id)
  AND (sqlc.narg(type)::TEXT IS NULL OR type = sqlc.narg(type)::TEXT)
  AND (sqlc.narg(category)::TEXT IS NULL OR category = sqlc.narg(category)::TEXT)
  AND (sqlc.narg(start_date)::DATE IS NULL OR date >= sqlc.narg(start_date)::DATE)
  AND (sqlc.narg(end_date)::DATE IS NULL OR date <= sqlc.narg(end_date)::DATE)
  AND deleted_at IS NULL
ORDER BY date DESC;

-- name: ListAllRecords :many
SELECT *
FROM records
WHERE (sqlc.narg('user_id')::UUID IS NULL OR user_id = sqlc.narg('user_id')::UUID)
  AND (sqlc.narg('type')::TEXT IS NULL OR type = sqlc.narg('type')::TEXT)
  AND (sqlc.narg('category')::TEXT IS NULL OR category = sqlc.narg('category')::TEXT)
  AND (sqlc.narg('start_date')::DATE IS NULL OR date >= sqlc.narg('start_date')::DATE)
  AND (sqlc.narg('end_date')::DATE IS NULL OR date <= sqlc.narg('end_date')::DATE)
  AND deleted_at IS NULL
ORDER BY date DESC;

-- name: GetRecordByID :one
SELECT *
FROM records
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetViewerRecordByID :one
SELECT * 
FROM records
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: UpdateRecordByID :one
UPDATE records
SET
  amount = $3,
  type = $4,
  category = $5,
  note = $6,
  date = $7,
  updated_at = NOW()
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

-- name: PatchRecordByID :one
UPDATE records
SET
  amount = COALESCE(sqlc.narg(amount)::NUMERIC, amount),
  type = COALESCE(sqlc.narg(type)::TEXT, type),
  category = COALESCE(sqlc.narg(category)::TEXT, category),
  note = COALESCE(sqlc.narg(note)::TEXT, note),
  date = COALESCE(sqlc.narg(date)::TIMESTAMP, date),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteRecordByID :exec
UPDATE records
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: HardDeleteRecordByID :exec
DELETE FROM records
WHERE id = $1 AND user_id = $2;
