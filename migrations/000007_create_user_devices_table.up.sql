DO $$
BEGIN
    -- Create table if it doesn't exist
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public'
          AND table_name = 'user_devices'
    ) THEN
        CREATE TABLE public.user_devices (
            id text PRIMARY KEY,
            user_id text NULL,
            device_id varchar(150) NULL,
            platform varchar(20) NOT NULL,      -- 'ios' | 'android'
            provider varchar(20) NOT NULL,      -- 'fcm' | 'apns'
            token text NOT NULL,
            app_version varchar(50) NULL,
            device_model varchar(100) NULL,
            last_seen_at timestamptz NULL,
            revoked_at timestamptz NULL,
            is_active boolean DEFAULT true,
            is_deleted boolean DEFAULT false,
            created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
            updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,

            CONSTRAINT fk_user_devices_user
                FOREIGN KEY (user_id) REFERENCES public.users(id)
                ON DELETE CASCADE
        );
    END IF;

    /*
      If the table already existed from previous runs, you might already have
      a non-partial unique index on (provider, token). That index will cause
      conflicts when re-registering the same token after soft-unregister.
      So we drop it if it exists, then re-create as partial unique.
    */
    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname = 'idx_user_devices_provider_token_unique'
    ) THEN
        DROP INDEX public.idx_user_devices_provider_token_unique;
    END IF;

    -- Scenario #1: Unique token only among ACTIVE (not revoked) and not deleted devices
    CREATE UNIQUE INDEX IF NOT EXISTS idx_user_devices_provider_token_unique
        ON public.user_devices(provider, token)
        WHERE revoked_at IS NULL AND is_deleted = false;

    -- One active device row per user_id + device_id (recommended to also be partial)
    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname = 'idx_user_devices_user_device_unique'
    ) THEN
        DROP INDEX public.idx_user_devices_user_device_unique;
    END IF;

    CREATE UNIQUE INDEX IF NOT EXISTS idx_user_devices_user_device_unique
        ON public.user_devices(user_id, device_id)
        WHERE user_id IS NOT NULL
          AND device_id IS NOT NULL
          AND revoked_at IS NULL
          AND is_deleted = false;

    -- Helpful query index for sending notifications / listing active devices per user
    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname = 'idx_user_devices_user_active'
    ) THEN
        DROP INDEX public.idx_user_devices_user_active;
    END IF;

    CREATE INDEX IF NOT EXISTS idx_user_devices_user_active
        ON public.user_devices(user_id)
        WHERE revoked_at IS NULL AND is_deleted = false AND is_active = true;

END$$;
