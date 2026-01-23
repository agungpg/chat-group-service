DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public'
          AND table_name = 'friendships'
    ) THEN
        CREATE TABLE public.friendships (
            user_low_id  text NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
            user_high_id text NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
            created_at   timestamptz NOT NULL DEFAULT now(),

            PRIMARY KEY (user_low_id, user_high_id),
            CONSTRAINT friendships_no_self CHECK (user_low_id <> user_high_id)
        );
    END IF;
END$$;
