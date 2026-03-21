-- CreateEnum
CREATE TYPE "Role" AS ENUM ('admin', 'user');
CREATE TYPE "DeviceStatus" AS ENUM ('online', 'offline', 'in_game');
CREATE TYPE "GameStatus" AS ENUM ('lobby', 'running', 'finished');

-- Users
CREATE TABLE "users" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "username" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "role" "Role" NOT NULL DEFAULT 'user',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "users_username_key" ON "users"("username");

-- Admin sessions
CREATE TABLE "admin_sessions" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "user_id" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "expires_at" TIMESTAMP(3) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "admin_sessions_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "admin_sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX "admin_sessions_token_key" ON "admin_sessions"("token");

-- Devices
CREATE TABLE "devices" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "device_id" TEXT NOT NULL,
    "status" "DeviceStatus" NOT NULL DEFAULT 'offline',
    "last_seen" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "devices_pkey" PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "devices_device_id_key" ON "devices"("device_id");

-- Games
CREATE TABLE "games" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "code" TEXT NOT NULL,
    "status" "GameStatus" NOT NULL DEFAULT 'lobby',
    "settings" JSONB NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "started_at" TIMESTAMP(3),
    "ended_at" TIMESTAMP(3),
    CONSTRAINT "games_pkey" PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "games_code_key" ON "games"("code");

-- Teams
CREATE TABLE "teams" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "game_id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    CONSTRAINT "teams_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "teams_game_id_fkey" FOREIGN KEY ("game_id") REFERENCES "games"("id") ON DELETE CASCADE
);

-- Players
CREATE TABLE "players" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "game_id" TEXT NOT NULL,
    "team_id" TEXT,
    "device_id" TEXT NOT NULL,
    "nickname" TEXT NOT NULL,
    "score" INTEGER NOT NULL DEFAULT 0,
    "kills" INTEGER NOT NULL DEFAULT 0,
    "deaths" INTEGER NOT NULL DEFAULT 0,
    "is_alive" BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT "players_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "players_game_id_fkey" FOREIGN KEY ("game_id") REFERENCES "games"("id") ON DELETE CASCADE,
    CONSTRAINT "players_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "teams"("id"),
    CONSTRAINT "players_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices"("device_id")
);
CREATE UNIQUE INDEX "players_game_id_device_id_key" ON "players"("game_id", "device_id");

-- Game events
CREATE TABLE "game_events" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "game_id" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "payload" JSONB NOT NULL DEFAULT '{}',
    "timestamp" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "game_events_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "game_events_game_id_fkey" FOREIGN KEY ("game_id") REFERENCES "games"("id") ON DELETE CASCADE
);
