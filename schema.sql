create table currency
(
    id   smallint PRIMARY KEY,
    code text NOT NULL,
    CHECK (code <> ''),
    CHECK (upper(code) = code),
    UNIQUE (code)
);

create table account
(
    id          text PRIMARY KEY,
    currency_id smallint not null references currency (id) on delete restrict deferrable,
    balance     numeric,
    CHECK (balance >= 0),
    CHECK (id <> '')
) partition by hash (id);

create table payment
(
    id              bigserial PRIMARY KEY,
    time            timestamp with time zone,
    from_account_id text     not null references account (id) on delete restrict deferrable,
    to_account_id   text     not null references account (id) on delete restrict deferrable,
    currency_id     smallint not null references currency (id) on delete restrict deferrable,
    amount          numeric  not null,
    CHECK (amount > 0)
) partition by hash (id);

create table account_part0 partition of account for values with (modulus 10, remainder 0);
create table account_part1 partition of account for values with (modulus 10, remainder 1);
create table account_part2 partition of account for values with (modulus 10, remainder 2);
create table account_part3 partition of account for values with (modulus 10, remainder 3);
create table account_part4 partition of account for values with (modulus 10, remainder 4);
create table account_part5 partition of account for values with (modulus 10, remainder 5);
create table account_part6 partition of account for values with (modulus 10, remainder 6);
create table account_part7 partition of account for values with (modulus 10, remainder 7);
create table account_part8 partition of account for values with (modulus 10, remainder 8);
create table account_part9 partition of account for values with (modulus 10, remainder 9);

create table payment_part0 partition of payment for values with (modulus 10, remainder 0);
create table payment_part1 partition of payment for values with (modulus 10, remainder 1);
create table payment_part2 partition of payment for values with (modulus 10, remainder 2);
create table payment_part3 partition of payment for values with (modulus 10, remainder 3);
create table payment_part4 partition of payment for values with (modulus 10, remainder 4);
create table payment_part5 partition of payment for values with (modulus 10, remainder 5);
create table payment_part6 partition of payment for values with (modulus 10, remainder 6);
create table payment_part7 partition of payment for values with (modulus 10, remainder 7);
create table payment_part8 partition of payment for values with (modulus 10, remainder 8);
create table payment_part9 partition of payment for values with (modulus 10, remainder 9);

create index on payment using btree (from_account_id, time desc);
create index on payment using btree (to_account_id, time desc);
