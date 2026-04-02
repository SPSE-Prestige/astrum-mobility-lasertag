export interface PlayerSessionResponse {
  player: PlayerDTO;
  game: GameDTO;
  team: TeamDTO | null;
  remaining_time: number;
  leaderboard: LeaderboardPlayerDTO[];
  events: GameEventDTO[];
}

export interface PlayerDTO {
  nickname: string;
  score: number;
  kills: number;
  deaths: number;
  is_alive: boolean;
  kill_streak: number;
  weapon_level: number;
  shots_fired: number;
}

export interface GameDTO {
  code: string;
  status: "lobby" | "running" | "finished";
  settings: GameSettingsDTO;
}

export interface GameSettingsDTO {
  respawn_delay: number;
  game_duration: number;
  friendly_fire: boolean;
  max_players: number;
  score_per_kill: number;
  kills_per_upgrade: number;
}

export interface TeamDTO {
  name: string;
  color: string;
}

export interface LeaderboardPlayerDTO {
  nickname: string;
  score: number;
  kills: number;
  deaths: number;
  shots_fired: number;
  team_name?: string;
  team_color?: string;
  is_current: boolean;
}

export interface GameEventDTO {
  id: string;
  type: string;
  payload: Record<string, unknown>;
  timestamp: string;
}

export interface ApiError {
  error: {
    code: string;
    message: string;
  };
}
