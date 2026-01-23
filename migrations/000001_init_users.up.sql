DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users') THEN
        CREATE TABLE public.users (
            id text PRIMARY KEY,
            username varchar(100) NOT NULL,
            email varchar(50) NOT NULL UNIQUE,
            password varchar(255) NOT NULL,
            email_verified_at TIMESTAMPTZ NULL,
            is_deleted boolean DEFAULT false,
            created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
            updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END$$; 