create table if not exists conversations (
    id bigserial primary key,
    type conversation_type not null,
    created_at timestamptz not null default now()
);

create index idx_conversations_type on conversations(type);

