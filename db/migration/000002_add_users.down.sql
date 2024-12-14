ALTER table if exists "accounts" drop CONSTRAINT IF exists "owner_currency_key";

ALTER table if exists "accounts" drop CONSTRAINT IF exists "accounts_owner_fkey";

DROP TABLE if exists "users";