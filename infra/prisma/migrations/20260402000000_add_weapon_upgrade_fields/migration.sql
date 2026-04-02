-- AlterTable: Add weapon upgrade tracking columns to players
ALTER TABLE "players" ADD COLUMN "kill_streak" INTEGER NOT NULL DEFAULT 0;
ALTER TABLE "players" ADD COLUMN "weapon_level" INTEGER NOT NULL DEFAULT 0;
