create extension if not exists "pg_uuidv7";

create table if not exists messages (
    id uuid primary key default uuidv7(),
    conversation_id bigint not null,
    sender_id bigint not null,
    content text not null,
    sent_at timestamptz not null default now(),
    
    constraint fk_messages_conversation foreign key (conversation_id) references conversations(id) on delete cascade,
    constraint fk_messages_sender foreign key (sender_id) references users(user_id) on delete cascade
);

create index idx_messages_conversation_sender_with_id on messages(conversation_id, sender_id) include (id);
create index idx_messages_sender_id on messages(sender_id);
create index idx_messages_conversation on messages(conversation_id, id desc);