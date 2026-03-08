create type oauth_login_result as (
    user_id bigint,
    session_id uuid
);

create or replace function oauth_login(
    p_provider text,
    p_provider_user_id text,
    p_display_name text,
    p_device_id uuid
)
returns oauth_login_result
language plpgsql
as $$
declare
    v_user_id bigint;
    v_session_id uuid;
begin
    -- 1. try to find existing identity
    select ai.user_id
    into v_user_id
    from auth_identities ai
    where ai.provider = p_provider
      and ai.provider_user_id = p_provider_user_id;

    -- 2. if not found → create new user + identity 
    if v_user_id is null then
        begin
            -- create new user
            insert into users (display_name)
            values (p_display_name)
            returning users.user_id into v_user_id;

            -- create auth identity
            insert into auth_identities (user_id, provider, provider_user_id)
            values (v_user_id, p_provider, p_provider_user_id);

        exception
            when unique_violation then
                -- if have race condition: another request created the same identity first
                -- get user_id again from auth_identities
                select ai.user_id
                into v_user_id
                from auth_identities ai
                where ai.provider = p_provider
                  and ai.provider_user_id = p_provider_user_id;
                
                -- if still not found → raise exception
                if v_user_id is null then
                    raise exception 'Race condition: could not find user after unique violation';
                end if;
        end;
    end if;

    -- 3. check if user is active
    if not exists (select 1 from users where user_id = v_user_id and is_active = true) then
        raise exception 'User account is deactivated or not found';
    end if;

    -- 4. create session
    insert into sessions (user_id, device_id)
    values (v_user_id, p_device_id)
    returning sessions.session_id into v_session_id;

    return (v_user_id, v_session_id);

end;
$$;