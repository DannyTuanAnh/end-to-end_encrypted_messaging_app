create table if not exists message_reads (
    message_id bigint not null,
    user_id bigint not null,
    read_at timestamptz not null default now(),

    primary key (message_id, user_id),

    constraint fk_message_reads_message foreign key (message_id) references messages(id) on delete cascade,
    constraint fk_message_reads_user foreign key (user_id) references users(user_id) on delete cascade
);

create index idx_message_reads_user_id on message_reads(user_id);