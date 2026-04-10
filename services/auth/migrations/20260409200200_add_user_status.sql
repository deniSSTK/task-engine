-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "status" character varying NOT NULL DEFAULT 'ACTIVE';
