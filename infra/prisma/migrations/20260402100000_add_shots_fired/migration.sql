-- AlterTable: add shots_fired column to players for accuracy tracking
ALTER TABLE "players" ADD COLUMN "shots_fired" INTEGER NOT NULL DEFAULT 0;
