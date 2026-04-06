-- +goose up

CREATE INDEX idx_records_user_date_active
ON records(user_id, date DESC)
WHERE deleted_at IS NULL;

-- income vs expense aggregation
CREATE INDEX idx_records_user_type_active ON records(user_id, type) WHERE deleted_at IS NULL;


-- +goose down

DROP INDEX IF EXISTS idx_records_user_date_active;
DROP INDEX IF EXISTS idx_records_user_type_active;
