-- +goose up

-- 1. Drop existing foreign key constraint
ALTER TABLE records
  DROP CONSTRAINT IF EXISTS records_user_id_fkey;

-- 2. Re-add foreign key WITHOUT CASCADE
ALTER TABLE records
  ADD CONSTRAINT records_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES users(id);

-- 3. Add created_by column
ALTER TABLE records
  ADD COLUMN created_by UUID;

-- 4. Add foreign key for created_by
ALTER TABLE records
  ADD CONSTRAINT records_created_by_fkey
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;


-- +goose down

-- 1. Drop created_by constraint
ALTER TABLE records
  DROP CONSTRAINT IF EXISTS records_created_by_fkey;

-- 2. Drop created_by column
ALTER TABLE records
  DROP COLUMN created_by;

-- 3. Drop modified user_id constraint
ALTER TABLE records
  DROP CONSTRAINT IF EXISTS records_user_id_fkey;

-- 4. Restore original constraint with CASCADE
ALTER TABLE records
  ADD CONSTRAINT records_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
