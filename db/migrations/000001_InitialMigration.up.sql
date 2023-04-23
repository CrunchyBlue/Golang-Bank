create table account
(
    id         bigserial
        primary key,
    owner      varchar                                not null,
    balance    bigint                                 not null,
    currency   varchar                                not null,
    created_at timestamp with time zone default now() not null
);

alter table account
    owner to root;

create index accounts_owner_idx
    on account (owner);

create table entry
(
    id         bigserial
        primary key,
    account_id bigint                                 not null
        references account,
    amount     bigint                                 not null,
    created_at timestamp with time zone default now() not null
);

comment on column entry.amount is 'Can be negative or positive';

alter table entry
    owner to root;

create index entries_account_id_idx
    on entry (account_id);

create table transfer
(
    id                     bigserial
        primary key,
    source_account_id      bigint                                 not null
        references account,
    destination_account_id bigint                                 not null
        references account,
    amount                 bigint                                 not null,
    created_at             timestamp with time zone default now() not null
);

comment on column transfer.amount is 'Must be positive';

alter table transfer
    owner to root;

create index transfers_source_account_id_idx
    on transfer (source_account_id);

create index transfers_destination_account_id_idx
    on transfer (destination_account_id);

create index transfers_source_account_id_destination_account_id_idx
    on transfer (source_account_id, destination_account_id);
