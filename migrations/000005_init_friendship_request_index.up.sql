DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname  = 'idx_friend_requests_addressee_status'
    ) THEN
        CREATE INDEX idx_friend_requests_addressee_status
            ON public.friend_requests (addressee_id, status);
    END IF;
END$$;