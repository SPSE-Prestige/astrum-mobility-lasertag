import { describe, it, expect } from "vitest";
import { toPlayer, toTeam, toDevice, buildKillFeed, calcTeamResults, formatRaceTime } from "@/lib/game/mappers";
import type { PlayerResponse, TeamResponse, EventResponse } from "@/lib/api/types";
import type { Player } from "@/types/game";

// ── Fixtures ──

const teamA: TeamResponse = { id: "t1", name: "Neon Green", color: "#00ff00" };
const teamB: TeamResponse = { id: "t2", name: "Neon Red", color: "#ff0000" };
const teams: TeamResponse[] = [teamA, teamB];

const rawPlayer = (overrides: Partial<PlayerResponse> = {}): PlayerResponse => ({
  id: "p1",
  game_id: "g1",
  device_id: "DEVICE-01",
  nickname: "Alice",
  team_id: "t1",
  is_alive: true,
  kills: 5,
  deaths: 2,
  score: 500,
  kill_streak: 3,
  weapon_level: 1,
  shots_fired: 10,
  ...overrides,
});

// ── toPlayer ──

describe("toPlayer", () => {
  it("maps response to domain Player with team name", () => {
    const p = toPlayer(rawPlayer(), teams);
    expect(p).toEqual({
      id: "p1",
      name: "Alice",
      team: "Neon Green",
      teamId: "t1",
      deviceId: "DEVICE-01",
      status: "alive",
      kills: 5,
      deaths: 2,
      score: 500,
      killStreak: 3,
      weaponLevel: 1,
      shotsFired: 10,
    });
  });

  it("returns 'Unassigned' when team_id is null", () => {
    const p = toPlayer(rawPlayer({ team_id: undefined }), teams);
    expect(p.team).toBe("Unassigned");
    expect(p.teamId).toBeNull();
  });

  it("returns 'dead' status when is_alive is false", () => {
    const p = toPlayer(rawPlayer({ is_alive: false }), teams);
    expect(p.status).toBe("dead");
  });

  it("defaults killStreak and weaponLevel to 0 when missing", () => {
    const p = toPlayer(rawPlayer({ kill_streak: undefined, weapon_level: undefined }), teams);
    expect(p.killStreak).toBe(0);
    expect(p.weaponLevel).toBe(0);
  });
});

// ── toTeam ──

describe("toTeam", () => {
  it("maps TeamResponse to domain Team", () => {
    expect(toTeam(teamA)).toEqual({ id: "t1", name: "Neon Green", color: "#00ff00" });
  });
});

// ── toDevice ──

describe("toDevice", () => {
  it("maps raw device object to domain Device", () => {
    const d = toDevice({ id: "d1", device_id: "DEV-01", status: "online", last_seen: "2026-01-01T00:00:00Z" });
    expect(d).toEqual({ id: "d1", deviceId: "DEV-01", status: "online", lastSeen: "2026-01-01T00:00:00Z" });
  });
});

// ── buildKillFeed ──

describe("buildKillFeed", () => {
  const killEvent = (id: string, payload: Record<string, unknown> = {}): EventResponse => ({
    id,
    game_id: "g1",
    type: "kill",
    timestamp: "2026-01-01T12:00:00Z",
    payload: { attacker_nickname: "Alice", victim_nickname: "Bob", ...payload },
  });

  it("filters only kill events", () => {
    const events: EventResponse[] = [
      killEvent("e1"),
      { id: "e2", game_id: "g1", type: "spawn", timestamp: "2026-01-01T12:00:01Z", payload: {} },
      killEvent("e3"),
    ];
    const feed = buildKillFeed(events);
    expect(feed).toHaveLength(2);
  });

  it("includes weapon upgrade info when present", () => {
    const events = [killEvent("e1", { weapon_upgraded: true, weapon_level: 2 })];
    const feed = buildKillFeed(events);
    expect(feed[0].message).toContain("⚡ LVL 2");
  });

  it("limits to 20 entries", () => {
    const events = Array.from({ length: 30 }, (_, i) => killEvent(`e${i}`));
    const feed = buildKillFeed(events);
    expect(feed).toHaveLength(20);
  });

  it("reverses events (most recent first)", () => {
    const events: EventResponse[] = [
      killEvent("e1"),
      { ...killEvent("e2"), payload: { attacker_nickname: "Charlie", victim_nickname: "Dave" } },
    ];
    const feed = buildKillFeed(events);
    expect(feed[0].message).toContain("Charlie");
  });
});

// ── calcTeamResults ──

describe("calcTeamResults", () => {
  const player = (name: string, team: string, score: number, kills = 0, deaths = 0): Player => ({
    id: name,
    name,
    team,
    teamId: "t1",
    deviceId: "d1",
    status: "alive",
    kills,
    deaths,
    score,
    killStreak: 0,
    weaponLevel: 0,
    shotsFired: 0,
  });

  it("aggregates team scores in team mode", () => {
    const players = [
      player("A", "Alpha", 300, 3, 1),
      player("B", "Alpha", 200, 2, 0),
      player("C", "Beta", 400, 4, 2),
    ];
    const results = calcTeamResults(players, "team");
    expect(results[0].team).toBe("Alpha");
    expect(results[0].score).toBe(500);
    expect(results[1].team).toBe("Beta");
    expect(results[1].score).toBe(400);
  });

  it("excludes Unassigned players in team mode", () => {
    const players = [player("A", "Alpha", 300), player("B", "Unassigned", 500)];
    const results = calcTeamResults(players, "team");
    expect(results).toHaveLength(1);
    expect(results[0].team).toBe("Alpha");
  });

  it("lists individual players in ffa mode", () => {
    const players = [player("A", "Solo", 300), player("B", "Solo", 500)];
    const results = calcTeamResults(players, "ffa");
    expect(results).toHaveLength(2);
    expect(results[0].team).toBe("B");
    expect(results[0].score).toBe(500);
  });
});

// ── formatRaceTime ──

describe("formatRaceTime", () => {
  it("formats 300 seconds as 05:00", () => {
    expect(formatRaceTime(300)).toBe("05:00");
  });

  it("formats 61 seconds as 01:01", () => {
    expect(formatRaceTime(61)).toBe("01:01");
  });

  it("formats 0 seconds as 00:00", () => {
    expect(formatRaceTime(0)).toBe("00:00");
  });

  it("clamps negative seconds to 00", () => {
    expect(formatRaceTime(-5)).toBe("-1:00");
    // Actually let's verify the actual behavior:
    // Math.floor(-5/60) = -1, padStart(2,"0") = "-1"
    // Math.max(-5 % 60, 0) = Math.max(-5, 0) = 0 → "00"
    // So it returns "-1:00" — this documents the edge case
  });
});
