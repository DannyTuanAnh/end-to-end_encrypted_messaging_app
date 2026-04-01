-- name: GetProfileByUserId :one
SELECT p.*, u.uuid
FROM profiles p
join users u on p.user_id = u.user_id
WHERE p.user_id = $1;

-- name: GetProfilesByUserUUID :one
SELECT p.name, p.user_id, p.avatar_url
FROM profiles p
JOIN users u ON p.user_id = u.user_id
WHERE u.uuid = $1 AND u.is_active = true;

-- name: CreateProfile :one
INSERT INTO profiles (user_id, name, email, birthday, avatar_url) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateProfileByUserId :one
UPDATE profiles SET 
    name = COALESCE(sqlc.narg('name'), name),
    birthday = COALESCE(sqlc.narg('birthday'), birthday),
    email = COALESCE(sqlc.narg('email'), email),
    phone = COALESCE(sqlc.narg('phone'), phone),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = now()
WHERE user_id = $1 
RETURNING *;

