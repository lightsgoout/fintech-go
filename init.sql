create type currency as enum (
    'USD',
    'EUR',
    'RUB'
    );

create table account
(
    id       text PRIMARY KEY,
    currency currency not null,
    balance  numeric,
    CHECK (balance >= 0),
    CHECK (id <> '')
);

create table payment
(
    id              bigserial PRIMARY KEY,
    time            timestamp with time zone,
    from_account_id text     not null references account (id) on delete restrict deferrable,
    to_account_id   text     not null references account (id) on delete restrict deferrable,
    currency        currency not null,
    amount          numeric  not null,
    CHECK (amount > 0)
);

create index on payment using btree (from_account_id, time desc);
create index on payment using btree (to_account_id, time desc);
create index on account using hash(currency);
