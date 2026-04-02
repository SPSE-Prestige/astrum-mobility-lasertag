import type { Device, GameMode, KillFeedEvent, Player, Team, TeamResult } from "@/types/game";
import type { EventResponse, PlayerResponse, TeamResponse } from "@/lib/api/types";

export function toPlayer(p: PlayerResponse, teams: TeamResponse[]): Player {
  const team = teams.find((t) => t.id === p.team_id);
  return {
    id: p.id,
    name: p.nickname,
    team: team?.name ?? "Unassigned",
    teamId: p.team_id ?? null,
    deviceId: p.device_id,
    status: p.is_alive ? "alive" : "dead",
    kills: p.kills,
    deaths: p.deaths,
    score: p.score,
    killStreak: p.kill_streak ?? 0,
    weaponLevel: p.weapon_level ?? 0,
  };
}

export function toTeam(t: TeamResponse): Team {
  return { id: t.id, name: t.name, color: t.color };
}

export function toDevice(d: { id: string; device_id: string; status: string; last_seen: string }): Device {
  return { id: d.id, deviceId: d.device_id, status: d.status, lastSeen: d.last_seen };
}

export function buildKillFeed(events: EventResponse[]): KillFeedEvent[] {
  return events
    .filter((e) => e.type === "kill")
    .reverse()
    .slice(0, 20)
    .map((e) => {
      const attacker = (e.payload?.attacker_nickname as string) ?? "?";
      const victim = (e.payload?.victim_nickname as string) ?? "?";
      const weaponUpgraded = e.payload?.weapon_upgraded === true;
      const weaponLevel = (e.payload?.weapon_level as number) ?? 0;

      let message = `${attacker} → ${victim}`;
      if (weaponUpgraded) message += ` ⚡ LVL ${weaponLevel}`;

      return {
        id: e.id,
        timestamp: new Date(e.timestamp).toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        }),
        message,
      };
    });
}

export function calcTeamResults(players: Player[], mode: GameMode): TeamResult[] {
  if (mode === "ffa") {
    return players
      .map((p) => ({ team: p.name, score: p.score, kills: p.kills, deaths: p.deaths }))
      .sort((a, b) => b.score - a.score);
  }
  const grouped = players
    .filter((p) => p.team !== "Unassigned")
    .reduce<Record<string, { score: number; kills: number; deaths: number }>>((acc, p) => {
      if (!acc[p.team]) acc[p.team] = { score: 0, kills: 0, deaths: 0 };
      acc[p.team].score += p.score;
      acc[p.team].kills += p.kills;
      acc[p.team].deaths += p.deaths;
      return acc;
    }, {});
  return Object.entries(grouped)
    .map(([team, s]) => ({ team, ...s }))
    .sort((a, b) => b.score - a.score);
}

export function formatRaceTime(seconds: number): string {
  const m = Math.floor(seconds / 60).toString().padStart(2, "0");
  const s = Math.max(seconds % 60, 0).toString().padStart(2, "0");
  return `${m}:${s}`;
}
