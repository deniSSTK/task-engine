-- Create "user_sessions" table
CREATE TABLE "user_sessions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "refresh_token" character varying NOT NULL,
  "ip" character varying NOT NULL,
  "user_agent" character varying NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "user_sessions_users_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "usersession_deleted_at" to table: "user_sessions"
CREATE INDEX "usersession_deleted_at" ON "user_sessions" ("deleted_at");
