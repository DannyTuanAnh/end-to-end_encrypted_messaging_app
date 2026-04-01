create table if not exists friend_requests (
    request_id bigserial primary key,
    sender_id bigint not null,
    receiver_id bigint not null,
    status friend_request_status not null default 'pending',
    send_at timestamptz not null default now(),

    constraint fk_friend_requests_sender foreign key (sender_id) references users(user_id) on delete cascade,
    constraint fk_friend_requests_receiver foreign key (receiver_id) references users(user_id) on delete cascade,
    constraint check_no_self_request check (sender_id <> receiver_id)
);

create unique index unique_friend_request on friend_requests(least(sender_id, receiver_id), greatest(sender_id, receiver_id));
create index idx_friend_requests_sender_receiver on friend_requests(sender_id, receiver_id);
create index idx_friend_requests_receiver_id on friend_requests(receiver_id);


