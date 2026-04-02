-- AlterTable: add session_code for mobile app player authentication
ALTER TABLE "players" ADD COLUMN "session_code" VARCHAR(6);

-- Create unique index (allows nulls for backwards compatibility)
CREATE UNIQUE INDEX "players_session_code_key" ON "players" ("session_code") WHERE "session_code" IS NOT NULL;
