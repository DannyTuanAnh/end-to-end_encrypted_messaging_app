-- name: AddFriendById :one
select * from send_friend_request($1, $2);

-- name: AcceptFriendRequestById :one
select * from accept_friend_request($1, $2);

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
    select user1_id as id from friendships where user2_id = $1
    union
    select user2_id as id from friendships where user1_id = $1
)

select 
    u.uuid, 
    coalesce(p.name, u.display_name) as name, 
    p.avatar_url,
    u.is_active
from users u
left join profiles p on u.user_id = p.user_id
join friend_ids f on u.user_id = f.id
order by coalesce(p.name, u.display_name);

-- name: RejectFriendRequestById :exec
delete from friend_requests
where request_id = $1 and receiver_id = $2 and status = 'pending';

-- name: RemoveFriendById :one
select remove_friend_by_uuid($1, $2) as deleted_count;