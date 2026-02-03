BEGIN;

DROP TABLE IF EXISTS files;

-- Drop enum if it exists
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_scope_enum') THEN
    DROP TYPE file_scope_enum;
  END IF;
END
$$;

COMMIT;