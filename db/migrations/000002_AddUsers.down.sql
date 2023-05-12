alter table if exists account
    drop constraint if exists owner_currency_key;

alter table if exists account
    drop constraint if exists account_owner_fkey;

drop table if exists "user";