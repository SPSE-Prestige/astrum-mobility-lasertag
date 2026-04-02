// ── Backend DTO types (match backend/internal/delivery/http/dto.go) ──

export interface LoginResponse {
  token: string;
  expires_at: string;
}

export interface DeviceResponse {
  id: string;
  device_id: string;
  status: string;
  last_seen: string;
}

export interface GameSettingsDTO {
  respawn_delay: number;
  game_duration: number;
  friendly_fire: boolean;
  max_players: number;
  score_per_kill: number;
  kills_per_upgrade: number;
}

export interface GameResponse {
  id: string;
  code: string;
  status: string;
  settings: GameSettingsDTO;
  created_at: string;
  started_at?: string;
  ended_at?: string;
}

export interface TeamResponse {
  id: number;
  game_id: string;
  name: string;
  color: string;
}

export interface PlayerResponse {
  id: string;
  game_id: string;
  team_id?: number;
  device_id: string;
  nickname: string;
  score: number;
  kills: number;
  deaths: number;
  is_alive: boolean;
  kill_streak: number;
  weapon_level: number;
  shots_fired: number;
  session_code?: string;
}

export interface GameFullResponse {
  game: GameResponse;
  teams: TeamResponse[];
  players: PlayerResponse[];
  events: EventResponse[];
}

export interface EventResponse {
  id: string;
  game_id: string;
  type: string;
  payload: Record<string, unknown>;
  timestamp: string;
}
