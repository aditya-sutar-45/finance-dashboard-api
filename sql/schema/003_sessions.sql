-- +goose up
CREATE TABLE sessions (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  user_email TEXT UNIQUE NOT NULL,
  refresh_token VARCHAR(512) NOT NULL,
  is_revoked BOOL NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP
);

-- +goose down
DROP TABLE sessions;
