BEGIN;

-- Create enum for scope (if it doesn't exist)
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_scope_enum') THEN
    CREATE TYPE file_scope_enum AS ENUM (
      'avatar',
      'chat_attachment',
      'task_attachment'
    );
  END IF;
END
$$;
-- Create enum for status (if it doesn't exist)
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_status_enum') THEN
    CREATE TYPE file_status_enum AS ENUM (
      'pending',
      'ready',
      'deleted'
    );
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS files (
  id                uuid PRIMARY KEY,

  created_by        text NOT NULL,
  bucket            text NOT NULL,
  object_key        text NOT NULL,

  original_filename text,
  content_type      text NOT NULL,
  size_bytes        bigint NOT NULL CHECK (size_bytes >= 0),

  scope             file_scope_enum NOT NULL,
  
  status             file_status_enum NOT NULL,

  created_at        timestamptz NOT NULL DEFAULT now(),
  deleted_at        timestamptz,

  CONSTRAINT files_created_by_fk
    FOREIGN KEY (created_by) REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

  -- Prevent duplicates across bucket+key
  CONSTRAINT files_bucket_object_key_uk
    UNIQUE (bucket, object_key)
);

-- Helpful indexes
CREATE INDEX IF NOT EXISTS idx_files_created_by ON files(created_by);
CREATE INDEX IF NOT EXISTS idx_files_scope ON files(scope);
CREATE INDEX IF NOT EXISTS idx_files_created_at ON files(created_at DESC);

-- Fast queries for "active" (not deleted) files
CREATE INDEX IF NOT EXISTS idx_files_not_deleted ON files(id)
WHERE deleted_at IS NULL;

COMMIT;
