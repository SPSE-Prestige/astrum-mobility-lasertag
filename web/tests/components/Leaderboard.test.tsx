import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Leaderboard } from "@/components/race-control/Leaderboard";
import type { Player } from "@/types/game";

const mockPlayer = (name: string, kills: number, deaths: number, score: number, weaponLevel = 0, shotsFired = 0): Player => ({
  id: name.toLowerCase(),
  name,
  team: "Alpha",
  teamId: "t1",
  deviceId: "d1",
  status: "alive",
  kills,
  deaths,
  score,
  killStreak: 0,
  weaponLevel,
  shotsFired,
});

describe("Leaderboard", () => {
  it("renders player rows", () => {
    const players = [mockPlayer("Alice", 5, 2, 500), mockPlayer("Bob", 3, 1, 300)];
    render(<Leaderboard players={players} gameMode="team" language="en" />);
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("Bob")).toBeInTheDocument();
  });

  it("shows weapon level when > 0", () => {
    const players = [mockPlayer("Alice", 5, 2, 500, 2)];
    render(<Leaderboard players={players} gameMode="team" language="en" />);
    expect(screen.getByText("LVL 2")).toBeInTheDocument();
  });

  it("shows dash when weapon level is 0", () => {
    const players = [mockPlayer("Alice", 5, 2, 500, 0)];
    render(<Leaderboard players={players} gameMode="team" language="en" />);
    const dashes = screen.getAllByText("—");
    expect(dashes.length).toBeGreaterThanOrEqual(1);
  });

  it("shows Solo in ffa mode", () => {
    const players = [mockPlayer("Alice", 5, 2, 500)];
    render(<Leaderboard players={players} gameMode="ffa" language="en" />);
    expect(screen.getByText("Solo")).toBeInTheDocument();
  });

  it("renders Czech header", () => {
    render(<Leaderboard players={[]} gameMode="team" language="cs" />);
    expect(screen.getByText("Průběžné pořadí")).toBeInTheDocument();
  });

  it("renders English header", () => {
    render(<Leaderboard players={[]} gameMode="team" language="en" />);
    expect(screen.getByText("Live Leaderboard")).toBeInTheDocument();
  });

  it("shows accuracy percentage when shots fired > 0", () => {
    const players = [mockPlayer("Alice", 5, 2, 500, 0, 20)];
    render(<Leaderboard players={players} gameMode="team" language="en" />);
    expect(screen.getByText("25.0%")).toBeInTheDocument();
  });

  it("shows dash for accuracy when no shots fired", () => {
    const players = [mockPlayer("Alice", 0, 0, 0, 0, 0)];
    render(<Leaderboard players={players} gameMode="team" language="en" />);
    const dashes = screen.getAllByText("—");
    expect(dashes.length).toBeGreaterThanOrEqual(2); // accuracy dash + weapon dash
  });
});
