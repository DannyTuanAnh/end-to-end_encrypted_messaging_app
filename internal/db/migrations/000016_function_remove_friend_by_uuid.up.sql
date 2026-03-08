create or replace function remove_friend_by_uuid(
    p_current_user_id bigint,
    p_friend_uuid uuid
)
returns bigint  -- return number of deleted rows (0 or 1)
language plpgsql
as $$
declare
    v_friend_id bigint;
    v_deleted_count bigint;
begin
    -- 1. check friend user_id from uuid
    select user_id into v_friend_id
    from users
    where uuid = p_friend_uuid;
    
    -- 2. if not found → raise exception
    if v_friend_id is null then
        raise exception 'User with UUID % not found', p_friend_uuid;
    end if;
    
    -- 3. check if trying to remove yourself
    if v_friend_id = p_current_user_id then
        raise exception 'Cannot remove yourself as friend';
    end if;
    
    -- 4. Delete friendship
    delete from friendships
    where user1_id = least(p_current_user_id, v_friend_id)
      and user2_id = greatest(p_current_user_id, v_friend_id);
    
    -- 5. get number of deleted rows
    get diagnostics v_deleted_count = row_count;
    
    -- 6. if no rows deleted → raise exception (friendship not found)
    if v_deleted_count = 0 then
        raise exception 'Friendship not found between users';
    end if;
    
    return v_deleted_count;
end;
$$;