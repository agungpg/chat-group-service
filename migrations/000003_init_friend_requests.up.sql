DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public'
          AND table_name = 'friend_requests'
    ) THEN
        CREATE TABLE public.friend_requests (
            id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
            requester_id text NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
            addressee_id text NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
            status       public.friend_request_status NOT NULL DEFAULT 'pending',
            created_at   timestamptz NOT NULL DEFAULT now(),
            updated_at   timestamptz NOT NULL DEFAULT now(),

            CONSTRAINT friend_requests_no_self CHECK (requester_id <> addressee_id),
            CONSTRAINT friend_requests_unique_pair UNIQUE (requester_id, addressee_id)
        );
    END IF;
END$$;
