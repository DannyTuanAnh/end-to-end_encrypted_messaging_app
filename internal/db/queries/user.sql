-- name: GetUserByUUID :one
-- Get user info with friendship/friend request status
-- $1: target_user_uuid (user being searched)
-- $2: current_user_id (user doing the search)
SELECT 
    u.user_id,
    u.uuid,
    u.display_name,
    
    -- Friend request status (if exists)
    CASE
        WHEN fr.status = 'pending' AND fr.sender_id = $2 THEN 'sent'
        WHEN fr.status = 'pending' AND fr.receiver_id = $2 THEN 'received'
        ELSE fr.status
    END as friend_request_direction,
    
    -- Friendship status (if exists)
    CASE 
        WHEN f.user1_id IS NOT NULL THEN true
        ELSE false
    END as is_friend
    
FROM users u

-- Check if there's a pending/accepted friend request
LEFT JOIN friend_requests fr 
    ON (fr.sender_id = $2 AND fr.receiver_id = u.user_id)
    OR (fr.sender_id = u.user_id AND fr.receiver_id = $2)

-- Check if already friends
LEFT JOIN friendships f
    ON (f.user1_id = LEAST($2, u.user_id) AND f.user2_id = GREATEST($2, u.user_id))

WHERE u.uuid = $1; 

-- name: GetUUIDByUserId :one
SELECT uuid FROM users WHERE user_id = $1;

-- name: CreateUser :one
INSERT INTO users (display_name) VALUES ($1) RETURNING *;
