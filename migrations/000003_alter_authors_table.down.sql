BEGIN;
ALTER TABLE authors
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deleted_at;
COMMIT;