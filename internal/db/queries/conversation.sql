-- name: GetAllConversations :many
with user_conversations as (
    select conversation_id
    from conversation_members
    where user_id = $1
),

latest_messages as (
    select distinct on (conversation_id)
        id,
        conversation_id,
        content as last_message,
        sent_at as last_message_time
    from messages
    order by conversation_id, sent_at desc
)

select
    c.id as conversation_id,
    c.type as conversation_type,
    
    -- Name and avatar based on conversation type
    case 
        when c.type = 'private' then
            (select p.name
             from conversation_members cm
             join profiles p on cm.user_id = p.user_id
             where cm.conversation_id = c.id 
               and cm.user_id <> $1
             limit 1)
        else
            (select g.name
             from groups g
             where g.conversation_id = c.id)
    end as conversation_name,
    
    case 
        when c.type = 'private' then
            (select p.avatar_url
             from conversation_members cm
             join profiles p on cm.user_id = p.user_id
             where cm.conversation_id = c.id 
               and cm.user_id <> $1
             limit 1)
        else
            (select g.avatar_url
             from groups g
             where g.conversation_id = c.id)
    end as avatar_url,
    
    -- Last message info
    lm.last_message,
    lm.last_message_time,

    (mr.message_id is not null) as is_read

from user_conversations uc
join conversations c on uc.conversation_id = c.id
left join latest_messages lm on c.id = lm.conversation_id
left join message_reads mr on lm.id = mr.message_id and mr.user_id = $1
order by coalesce(lm.last_message_time, c.created_at) desc;

-- name: GetMessagesByConversationId :many
select 
    m.id,
    m.sender_id,
    m.conversation_id,
    m.content,
    m.sent_at,
    coalesce(p.name, u.display_name) as sender_name,
    p.avatar_url as sender_avatar,
    case
        when m.sender_id = $2 then 'sent'
        else 'received'
    end as message_direction
from messages m
join users u on u.user_id = m.sender_id
left join profiles p on p.user_id = m.sender_id
where m.conversation_id = $1
order by m.sent_at desc;

-- name: MarkMessagesAsRead :exec
insert into message_reads (message_id, user_id, read_at)
select m.id, $2, now()
from messages m
left join message_reads mr on m.id = mr.message_id and mr.user_id = $2
where m.conversation_id = $1
  and m.sender_id <> $2
  and mr.message_id is null;

-- name: CreateMessage :one
insert into messages (sender_id, conversation_id, content)
values ($1, $2, $3)
returning id, sender_id, conversation_id, content, sent_at;

-- name: CreateGroupConversation


-- name: AddGroupMembers :exec
insert into conversation_members (conversation_id, user_id)
values ($1, $2);

-- name: RemoveGroupMembers :exec
delete from conversation_members
where conversation_id = $1 and user_id = $2;

-- name: LeaveConversation :exec
delete from conversation_members
where conversation_id = $1 and user_id = $2;
