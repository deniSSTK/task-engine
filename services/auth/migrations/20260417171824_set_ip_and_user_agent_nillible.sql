-- Modify "user_sessions" table
ALTER TABLE "user_sessions" ALTER COLUMN "ip" DROP NOT NULL, ALTER COLUMN "user_agent" DROP NOT NULL;
