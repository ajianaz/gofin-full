-- Add verified column to users table
-- Existing users are treated as verified (backward compatible)
ALTER TABLE users ADD COLUMN IF NOT EXISTS verified BOOLEAN;
UPDATE users SET verified = true WHERE verified IS NULL;
ALTER TABLE users ALTER COLUMN verified SET NOT NULL;
ALTER TABLE users ALTER COLUMN verified SET DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_users_verified ON users(verified);
