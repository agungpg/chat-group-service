DO $$
BEGIN
    -- Drop indexes (if they exist)
    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname = 'idx_user_devices_provider_token_unique'
    ) THEN
        DROP INDEX public.idx_user_devices_provider_token_unique;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname = 'idx_user_devices_user_device_unique'
    ) THEN
        DROP INDEX public.idx_user_devices_user_device_unique;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
          AND indexname = 'idx_user_devices_user_active'
    ) THEN
        DROP INDEX public.idx_user_devices_user_active;
    END IF;

    -- Drop table
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public'
          AND table_name = 'user_devices'
    ) THEN
        DROP TABLE public.user_devices;
    END IF;
END$$;
