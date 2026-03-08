-- name: GetProfileByUserId :one
SELECT * FROM profiles WHERE user_id = $1;

-- name: GetProfilesByUserUUID :one
SELECT p.*
FROM profiles p
JOIN users u ON p.user_id = u.user_id
WHERE u.uuid = $1;

-- name: CreateProfile :one
INSERT INTO profiles (user_id, name, birthday, avatar_url) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateProfileByUserId :one
UPDATE profiles SET 
    name = COALESCE(sqlc.narg('name'), name),
    birthday = COALESCE(sqlc.narg('birthday'), birthday),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = now()
WHERE user_id = $1 
RETURNING *;

