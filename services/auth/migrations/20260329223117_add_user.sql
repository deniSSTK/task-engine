-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "second_name" character varying NULL,
  "full_name" character varying NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  "last_logined_at" timestamptz NOT NULL,
  "role" character varying NOT NULL DEFAULT 'user',
  PRIMARY KEY ("id")
);
-- Create index "user_deleted_at" to table: "users"
CREATE INDEX "user_deleted_at" ON "users" ("deleted_at");
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
