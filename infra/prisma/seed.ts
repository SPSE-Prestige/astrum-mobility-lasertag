import { PrismaClient } from "@prisma/client";
import * as bcrypt from "bcryptjs";

const prisma = new PrismaClient();

async function main() {
  console.log("🌱 Seeding database...");

  // ===== Default Admin User =====
  const passwordHash = await bcrypt.hash("admin123", 10);
  const admin = await prisma.user.upsert({
    where: { username: "admin" },
    update: {},
    create: {
      username: "admin",
      passwordHash,
      role: "admin",
    },
  });
  console.log(`✅ Admin user created: ${admin.username}`);

  // ===== Sample Weapons =====
  const weapons = [
    {
      name: "Pistol",
      damage: 15,
      fireRateMs: 500,
      ammo: 12,
      reloadTimeMs: 2000,
      fireMode: "single",
      accuracySpread: 0.1,
    },
    {
      name: "Assault Rifle",
      damage: 10,
      fireRateMs: 100,
      ammo: 30,
      reloadTimeMs: 3000,
      fireMode: "auto",
      accuracySpread: 0.2,
    },
    {
      name: "Shotgun",
      damage: 40,
      fireRateMs: 1000,
      ammo: 6,
      reloadTimeMs: 4000,
      fireMode: "single",
      accuracySpread: 0.5,
    },
    {
      name: "Sniper Rifle",
      damage: 50,
      fireRateMs: 2000,
      ammo: 5,
      reloadTimeMs: 3500,
      fireMode: "single",
      accuracySpread: 0.02,
    },
    {
      name: "SMG",
      damage: 8,
      fireRateMs: 80,
      ammo: 40,
      reloadTimeMs: 2500,
      fireMode: "auto",
      accuracySpread: 0.3,
    },
  ];

  for (const w of weapons) {
    const weapon = await prisma.weapon.upsert({
      where: { name: w.name },
      update: w,
      create: w,
    });
    console.log(`✅ Weapon created: ${weapon.name}`);
  }

  // ===== Demo Game Config =====
  const demoGame = await prisma.game.create({
    data: {
      name: "Demo Team Deathmatch",
      status: "pending",
      configJson: {
        duration_seconds: 300,
        max_players: 10,
        team_count: 2,
        game_mode: "team_deathmatch",
        player: {
          max_hp: 100,
          lives: -1,
          respawn_delay_seconds: 5,
          friendly_fire: false,
        },
        scoring: {
          points_per_hit: 1,
          points_per_kill: 10,
          teamkill_penalty: 5,
          headshot_multiplier: 2.0,
        },
        feedback: {
          sound_enabled: true,
          vibration_enabled: true,
          led_enabled: true,
          intensity: 7,
        },
      },
      teams: {
        create: [
          { name: "Red Team", color: "#FF0000" },
          { name: "Blue Team", color: "#0000FF" },
        ],
      },
    },
  });
  console.log(`✅ Demo game created: ${demoGame.name} (${demoGame.id})`);

  console.log("🌱 Seeding complete!");
}

main()
  .catch((e) => {
    console.error("❌ Seeding failed:", e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
