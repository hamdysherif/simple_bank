CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE IF EXISTS "accounts" ADD COLUMN user_id bigint NOT NULL;
ALTER TABLE IF EXISTS "accounts" ADD CONSTRAINT "account_currency_unique" UNIQUE( user_id, currency);
ALTER TABLE IF EXISTS "accounts" ADD CONSTRAINT "fk_accounts_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id");
