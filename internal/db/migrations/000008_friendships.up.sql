create table if not exists friendships (
    user1_id bigint not null,
    user2_id bigint not null,
    established_at timestamptz not null default now(),

    primary key (user1_id, user2_id),

    constraint fk_friendships_user1 foreign key (user1_id) references users(user_id) on delete cascade,
    constraint fk_friendships_user2 foreign key (user2_id) references users(user_id) on delete cascade,
    constraint check_user_order check (user1_id < user2_id)
);

create index idx_friendships_user2_id on friendships(user2_id);
