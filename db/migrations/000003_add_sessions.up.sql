create table "session"
(
    id            uuid primary key,
    username      varchar                                not null,
    refresh_token varchar                                not null,
    user_agent    varchar                                not null,
    client_ip     varchar                                not null,
    is_blocked    boolean                  default false not null,
    expires_at    timestamp with time zone               not null,
    created_at    timestamp with time zone default now() not null
);

alter table session
    add foreign key (username) references "user" (username);