-- name: OAuthLogin :one
select * from oauth_login($1, $2, $3, $4);

-- name: CheckSession :one
select user_id, revoked, revoke_at from sessions where session_id = $1 and device_id = $2;

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
update api_keys set is_active = false, revoked_at = now() where is_active = true

-- name: ValidateAPIKey :one
select is_active from api_keys where key_hash = $1;