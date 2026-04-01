create extension if not exists "pgcrypto";

create table if not exists sessions (
    session_id uuid primary key default gen_random_uuid(),
    user_id bigint not null,
    revoked boolean not null default false,
    revoke_at timestamptz,
    device_id uuid not null,
    created_at timestamptz not null default now(),

    constraint fk_sessions_user foreign key (user_id) references users(user_id) on delete cascade
);

create index idx_sessions_user_id on sessions(user_id);
create index idx_sessions_revoked on sessions(revoked);