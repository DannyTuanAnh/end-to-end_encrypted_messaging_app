create type friend_request_result as (
    status boolean,
    message text
);

create or replace function send_friend_request(p_sender_id bigint, p_receiver_id bigint)
returns friend_request_result
language plpgsql
as $$
declare
    v_existing_sender bigint;
begin

    if p_sender_id = p_receiver_id then
        return (false, 'You cannot send a friend request to yourself');
    end if;

    -- 1. check friendship
    if exists (
        select 1
        from friendships
        where user1_id = least(p_sender_id, p_receiver_id) and user2_id = greatest(p_sender_id, p_receiver_id)
    ) then
        return (false, 'already_friends');
    end if;

    -- 2. create friend request (if not exists)
    insert into friend_requests (sender_id, receiver_id)
    values (p_sender_id, p_receiver_id)
    on conflict do nothing;

    -- if failed to insert → check if there's an existing pending request 
    if not found then

        select sender_id into v_existing_sender 
        from friend_requests 
        where sender_id = p_receiver_id and receiver_id = p_sender_id and status = 'pending'
        limit 1;
        if v_existing_sender = p_receiver_id then
            return (false, 'This person has already sent you a friend request');
        else
            return (false, 'Friend request already sent');
        end if;

    end if;

    return (true, 'Friend request sent successfully');

end;
$$;