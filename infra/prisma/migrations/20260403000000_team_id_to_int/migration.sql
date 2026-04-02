-- Convert teams.id from TEXT (UUID) to SERIAL (auto-increment integer)
-- and players.team_id from TEXT to INTEGER accordingly.

-- 1. Drop FK constraint on players.team_id
ALTER TABLE "players" DROP CONSTRAINT IF EXISTS "players_team_id_fkey";

-- 2. Add a temporary integer column to teams
ALTER TABLE "teams" ADD COLUMN "new_id" SERIAL;

-- 3. Update players.team_id to reference the new integer id
ALTER TABLE "players" ADD COLUMN "new_team_id" INTEGER;

UPDATE "players" p
SET "new_team_id" = t."new_id"
FROM "teams" t
WHERE p."team_id" = t."id";

-- 4. Drop old columns and rename new ones
ALTER TABLE "players" DROP COLUMN "team_id";
ALTER TABLE "players" RENAME COLUMN "new_team_id" TO "team_id";

ALTER TABLE "teams" DROP CONSTRAINT "teams_pkey";
ALTER TABLE "teams" DROP COLUMN "id";
ALTER TABLE "teams" RENAME COLUMN "new_id" TO "id";
ALTER TABLE "teams" ADD PRIMARY KEY ("id");

-- 5. Restore FK constraint
ALTER TABLE "players" ADD CONSTRAINT "players_team_id_fkey"
  FOREIGN KEY ("team_id") REFERENCES "teams"("id") ON DELETE SET NULL ON UPDATE CASCADE;
