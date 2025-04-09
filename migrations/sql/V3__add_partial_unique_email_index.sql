-- Remove the old unique constraint
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_indexes WHERE indexname = 'users_email_key'
  ) THEN
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
END IF;
END
$$;

-- Create a partial unique index on email for active (non-deleted) users
CREATE UNIQUE INDEX uniq_active_email
    ON users(email)
    WHERE is_deleted = false;
