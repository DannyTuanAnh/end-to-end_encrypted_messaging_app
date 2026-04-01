create or replace function create_group_conversation(
    p_creator_id bigint,
    p_group_name text,
    p_member_ids bigint[]  -- Array of user IDs to add (including creator)
)
returns table (
    conversation_id bigint,
    group_name text,
    invalid_user_ids bigint[]
)
language plpgsql
as $$
declare
    v_valid_users bigint[];
begin
    -- 1. Validate creator is active
    if not exists (select 1 from users where user_id = p_creator_id and is_active = true) then
        raise exception 'Creator user is not active or does not exist';
    end if;
    
    -- 2. Check which users are valid (active)
    select 
        array_agg(u.user_id) filter (where u.is_active = true),
        array_agg(m_id) filter (where u.user_id is null or u.is_active = false)
    into v_valid_users, invalid_user_ids
    from unnest(p_member_ids) as m_id
    left join users u on u.user_id = m_id;

    -- 3. Ensure creator is included in valid users (if not already)
    if v_valid_users is null or not (p_creator_id = any(v_valid_users)) then
        v_valid_users := coalesce(array_append(v_valid_users, p_creator_id), array[p_creator_id]);
    end if;
    
    -- 4. Check minimum group size (creator already included in v_valid_users)
    if coalesce(array_length(v_valid_users, 1), 0) < 3 then
        raise exception 'Group must have at least 3 members. Valid users: %, Invalid users: %', 
            coalesce(array_length(v_valid_users, 1), 0), 
            coalesce(array_length(invalid_user_ids, 1), 0);
    end if;
    
    -- 5. If group name not provided, generate from first 4 valid users
    if p_group_name is null or trim(p_group_name) = '' then
        select string_agg(u.display_name, ', ')
        into group_name
        from (select display_name from users where user_id = any(v_valid_users) limit 4) u;
    else
        group_name := p_group_name;
    end if;
    
    -- 6. Create conversation
    insert into conversations (type)
    values ('group')
    returning id into conversation_id;
    
    -- 7. Create group info
    insert into groups (conversation_id, name)
    values (v_conversation_id, v_group_name);
    
    -- 8. Add all members
    insert into conversation_members (conversation_id, user_id, role)
    select 
        v_conversation_id, 
        u_id, 
        case when u_id = p_creator_id then 'admin' else 'member' end
    from unnest(v_valid_users) as u_id
    on conflict (conversation_id, user_id) do nothing;
    
    return next;
end;
$$;
