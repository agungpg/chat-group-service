DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'friend_request_status'
          AND n.nspname = 'public'
    ) THEN
        CREATE TYPE public.friend_request_status AS ENUM (
            'pending',
            'accepted',
            'declined',
            'cancelled'
        );
    END IF;
END$$;
