-- name: GetDashboardSummary :one
SELECT
  COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::TEXT AS total_income,
  COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::TEXT AS total_expense
FROM records
WHERE deleted_at IS NULL
  AND (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID);

-- name: GetCategoryAnalysis :many
SELECT
  category,
  SUM(amount)::FLOAT AS total
FROM records
WHERE deleted_at IS NULL
  AND (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID)
GROUP BY category
ORDER BY total DESC;

-- name: GetTrends :many
SELECT
  TO_CHAR(date, 'YYYY-MM') AS month,
  SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END)::FLOAT AS income,
  SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END)::FLOAT AS expense
FROM records
WHERE deleted_at IS NULL
  AND (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID)
GROUP BY month
ORDER BY month DESC;

-- name: GetRecent :many
SELECT *
FROM records
WHERE deleted_at IS NULL
  AND (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID)
ORDER BY date DESC
LIMIT 10;
