drop function if exists oauth_login(
    p_provider text,
    p_provider_user_id text,
    p_display_name text,
    p_email varchar(255),
    p_device_id uuid
);