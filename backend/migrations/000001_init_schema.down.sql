-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_analyses_submission_id;
DROP INDEX IF EXISTS idx_submissions_created_at;
DROP INDEX IF EXISTS idx_submissions_status;
DROP INDEX IF EXISTS idx_submissions_user_id;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS analyses;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "pgcrypto";
