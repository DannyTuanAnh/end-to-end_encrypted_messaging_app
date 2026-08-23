-- name: FindExistingIdentity :one
select user_id, status, revoked_at from auth_identities 
where provider = $1 and provider_user_id = $2;

-- name: CreateIdentity :exec
insert into auth_identities (user_id, provider, provider_user_id, email)
values ($1, $2, $3, $4);

-- name: ActiveIdentity :execresult
update auth_identities 
set status = 'active', revoked_at = null 
where provider = $1 
and provider_user_id = $2
and revoked_at > now() - interval '30 days';

-- name: DisableIdentity :execresult
update auth_identities set status = 'revoked', revoked_at = now() where provider = $1 and user_id = $2;

-- name: CreateSession :one
insert into sessions (user_id) values ($1) returning session_id;

-- name: CheckSession :one
select user_id, revoked, revoke_at
from sessions 
where session_id = $1;

-- name: RevokeSession :exec
update sessions set revoked = true, revoke_at = now() where session_id = $1;

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