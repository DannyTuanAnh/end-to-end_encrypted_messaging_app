create table if not exists auth_identities (
    id bigserial primary key,
    user_id bigint not null,
    provider varchar(50) not null,
    provider_user_id varchar(255) not null,
    email varchar(255),
    created_at timestamptz not null default now(),
    
    constraint fk_auth_user foreign key (user_id) references users(user_id) on delete cascade,
    unique (provider, provider_user_id)
);

create index idx_auth_identities_user_id on auth_identities(user_id);