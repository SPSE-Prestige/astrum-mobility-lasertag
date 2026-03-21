import { PrismaClient } from "@prisma/client";
import { randomBytes, scryptSync } from "crypto";

const prisma = new PrismaClient();

function hashPassword(password: string): string {
  const salt = randomBytes(16).toString("hex");
  const hash = scryptSync(password, salt, 64).toString("hex");
  return `${salt}:${hash}`;
}

async function main() {
  console.log("Seeding database...");

  // Create default admin user
  const adminPassword = hashPassword("admin123");
  await prisma.user.upsert({
    where: { username: "admin" },
    update: {},
    create: {
      username: "admin",
      email: "admin@lasertag.local",
      passwordHash: adminPassword,
      role: "admin",
    },
  });
  console.log("Created default admin user (admin / admin123)");

  // Insert sample weapons
  const weapons = [
    { name: "Laser Pistol",  damage: 15, fireRate: 400,  ammo: 30,  reloadTime: 1500 },
    { name: "Laser Rifle",   damage: 25, fireRate: 600,  ammo: 20,  reloadTime: 2000 },
    { name: "Shotgun",       damage: 40, fireRate: 1000, ammo: 8,   reloadTime: 2500 },
    { name: "Sniper",        damage: 80, fireRate: 1500, ammo: 5,   reloadTime: 3000 },
    { name: "SMG",           damage: 10, fireRate: 200,  ammo: 50,  reloadTime: 1800 },
  ];

  for (const weapon of weapons) {
    await prisma.weapon.upsert({
      where: { name: weapon.name },
      update: {},
      create: weapon,
    });
  }
  console.log(`Inserted ${weapons.length} sample weapons`);

  console.log("Seeding complete.");
}

main()
  .catch((e) => {
    console.error("Seed failed:", e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
