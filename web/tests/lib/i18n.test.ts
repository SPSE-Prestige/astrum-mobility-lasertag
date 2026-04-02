import { describe, it, expect } from "vitest";
import { t } from "@/lib/i18n";

describe("i18n t()", () => {
  it("returns Czech translation", () => {
    expect(t("phase.setup", "cs")).toBe("Nastavení hry");
  });

  it("returns English translation", () => {
    expect(t("phase.setup", "en")).toBe("Game Setup");
  });

  it("returns different values for each language", () => {
    const cs = t("setup.createGame", "cs");
    const en = t("setup.createGame", "en");
    expect(cs).not.toBe(en);
    expect(cs).toBe("Vytvořit hru");
    expect(en).toBe("Create Game");
  });

  it("covers all phase keys", () => {
    const keys = ["phase.1", "phase.setup", "phase.players", "phase.live", "phase.results"] as const;
    for (const key of keys) {
      expect(t(key, "cs")).toBeTruthy();
      expect(t(key, "en")).toBeTruthy();
    }
  });

  it("covers auth keys", () => {
    expect(t("auth.login", "cs")).toBe("Přihlášení do dashboardu");
    expect(t("auth.logout", "en")).toBe("Logout");
  });
});
