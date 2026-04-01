-- name: OAuthLogin :one
select f.user_id::bigint, f.session_id::uuid, f.profile_exists::boolean from oauth_login($1, $2, $3, $4, $5) as f;

-- name: CheckSession :one
select s.user_id, s.revoked, s.revoke_at 
from sessions as s 
join users as u on s.user_id = u.id
where session_id = $1 and device_id = $2 and u.is_active = true;

-- name: RevokeSession :exec
update sessions set revoked = true, revoke_at = now() where session_id = $1 and device_id = $2;

-- name: RevokeAllSessions :exec
update sessions set revoked = true, revoke_at = now() where user_id = $1;

-- name: CleanupSessionTable :exec
delete from sessions where revoked = true and revoke_at < now() - interval '1 days';

-- manage apikeys
-- name: CreateAPIKey :exec
insert into api_keys (key_hash) values ($1);

-- name: RevokeAPIKeyByKey :exec
update api_keys set is_active = false, revoked_at = now() where key_hash = $1;

-- name: RevokeAllAPIKeys :exec
update api_keys set is_active = false, revoked_at = now() where is_active = true;

-- name: ValidateAPIKey :one
select is_active from api_keys where key_hash = $1;