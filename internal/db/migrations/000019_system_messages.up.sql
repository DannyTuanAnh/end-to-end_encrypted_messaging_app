create type system_event_type as enum (
    'member_joined',
    'member_left',
    'member_removed',
    'member_added',
    'group_created',
    'group_name_changed',
    'group_avatar_changed',
    'member_promoted',
    'member_demoted'
);

create table if not exists system_messages (
    id uuid primary key default uuidv7(),
    conversation_id bigint not null,
    event_type system_event_type not null,
    actor_id bigint,
    target_id bigint,
    metadata jsonb,
    created_at timestamptz not null default now(),
    
    constraint fk_system_messages_conversation 
        foreign key (conversation_id) references conversations(id) on delete cascade,
    constraint fk_system_messages_actor 
        foreign key (actor_id) references users(user_id) on delete set null,
    constraint fk_system_messages_target 
        foreign key (target_id) references users(user_id) on delete set null
);

create index idx_system_messages_conversation on system_messages(conversation_id, created_at desc);
create index idx_system_messages_event_type on system_messages(event_type);

-- Trigger để auto log member changes
create or replace function log_member_change()
returns trigger
language plpgsql
as $$
begin
    if TG_OP = 'INSERT' then
        insert into system_messages (conversation_id, event_type, target_id)
        values (NEW.conversation_id, 'member_joined', NEW.user_id);
        
    elsif TG_OP = 'DELETE' then
        insert into system_messages (conversation_id, event_type, target_id)
        values (OLD.conversation_id, 'member_left', OLD.user_id);
    end if;
    
    return NEW;
end;
$$;

create trigger trigger_log_member_change
after insert or delete on conversation_members
for each row
execute function log_member_change();