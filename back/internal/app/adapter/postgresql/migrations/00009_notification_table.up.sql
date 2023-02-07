create table notifications
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id    bigint
        constraint fk_notifications_user
            references users,
    read       boolean,
    message    text,
    level      bigint default 1
);

create index idx_notifications_deleted_at
    on notifications (deleted_at);

create index idx_notifications_read
    on notifications (read);
