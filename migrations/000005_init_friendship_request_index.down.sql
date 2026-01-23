DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname  = 'idx_friend_requests_addressee_status'
    ) THEN
        DROP INDEX public.idx_friend_requests_addressee_status;
    END IF;
END$$;
