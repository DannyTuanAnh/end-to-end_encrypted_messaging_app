create or replace function accept_friend_request(
    p_request_id bigint,
    p_receiver_id bigint
)
returns table(sender_id bigint, receiver_name text)
language plpgsql
as $$

declare 
    v_sender_id bigint;
    v_receiver_id bigint;
    v_updated_count int;
    v_conversation_id bigint;
    v_receiver_name text;

begin
    -- 1. Accept friend request
    update friend_requests
    set status = 'accepted'
    where request_id = p_request_id and receiver_id = p_receiver_id and status = 'pending'
    returning friend_requests.sender_id, friend_requests.receiver_id
    into v_sender_id, v_receiver_id;

    -- 2. If no rows updated → raise exception
    get diagnostics v_updated_count = row_count;
    if v_updated_count = 0 then
        raise exception 'Friend request not found or already processed';
    end if;

    -- 3. create friendship
    insert into friendships (user1_id, user2_id)
    values (least(v_sender_id, v_receiver_id), greatest(v_sender_id, v_receiver_id))
    on conflict do nothing;

    if not found then
        raise exception 'Friendship already exists';
    end if;

    -- 4. create conversation for the new friend
    insert into conversations(type)
    values('private')
    returning conversation_id into v_conversation_id;

    insert into conversation_members(conversation_id, user_id)
    values (v_conversation_id, v_sender_id), (v_conversation_id, v_receiver_id);

    -- 5. create system message for the new conversation
    insert into system_messages (conversation_id, event_type, content)
    values (v_conversation_id, 'chat_created', 'You are now friends! Start chatting.');

    -- 5. get receiver name for notification
    select name from profiles where user_id = v_receiver_id into v_receiver_name;

    -- 6. return sender_id and receiver_name for notification
    sender_id := v_sender_id;
    receiver_name := v_receiver_name;
    return next;
end;
$$;




    