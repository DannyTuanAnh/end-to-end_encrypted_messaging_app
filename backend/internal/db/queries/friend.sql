-- name: AddFriendById :one
select f.status::boolean, f.message::text from send_friend_request($1, $2) as f;

-- name: AcceptFriendRequestById :one
select a.sender_id::bigint , a.receiver_name::text from accept_friend_request($1, $2) as a;

-- name: GetPendingFriendRequests :many
select fr.request_id, u.uuid, p.name, p.avatar_url, fr.send_at
from friend_requests fr
join users u on fr.sender_id = u.user_id
left join profiles p on fr.sender_id = p.user_id
where fr.receiver_id = $1 and fr.status = 'pending'
order by fr.send_at desc;

-- name: GetSentFriendRequests :many
select fr.request_id, u.uuid, p.name, p.avatar_url, fr.send_at
from friend_requests fr
join users u on fr.receiver_id = u.user_id
left join profiles p on fr.receiver_id = p.user_id
where fr.sender_id = $1 and fr.status = 'pending'
order by fr.send_at desc;

-- name: GetFriendsList :many
with friend_ids as (
    select f.user1_id as id from friendships f where f.user2_id = $1
    union all
    select f.user2_id as id from friendships f where f.user1_id = $1
)

select 
    u.uuid, 
    u.user_id,
    coalesce(p.name, u.display_name) as name, 
    p.avatar_url,
    u.is_active
from users u
left join profiles p on u.user_id = p.user_id
join friend_ids f on u.user_id = f.id
order by coalesce(p.name, u.display_name);

-- name: SearchFriendByName :many
with friend_ids as (
    select f.user1_id as id from friendships f where f.user2_id = $1
    union all
    select f.user2_id as id from friendships f where f.user1_id = $1
)

select 
    u.uuid, 
    u.user_id,
    coalesce(p.name, u.display_name) as name, 
    p.avatar_url,
    u.is_active
from users u
left join profiles p on u.user_id = p.user_id
join friend_ids f on u.user_id = f.id
where coalesce(p.name, u.display_name) ilike '%' || $2 || '%' 
order by coalesce(p.name, u.display_name);

-- name: RejectFriendRequestById :exec
delete from friend_requests
where request_id = $1 and receiver_id = $2 and status = 'pending';

-- name: RemoveFriendById :one
select remove_friend_by_uuid($1, $2) as deleted_count;