export interface PlayerData {
  nickname: string;
  score: number;
  kills: number;
  deaths: number;
  is_alive: boolean;
  kill_streak: number;
  weapon_level: number;
  shots_fired: number;
}

export interface GameSettings {
  respawn_delay: number;
  game_duration: number;
  friendly_fire: boolean;
  max_players: number;
  score_per_kill: number;
  kills_per_upgrade: number;
}

export interface GameData {
  code: string;
  status: 'running' | 'lobby' | 'finished';
  settings: GameSettings;
}

export interface TeamData {
  name: string;
  color: string;
}

export interface PlayerSession {
  player: PlayerData;
  game: GameData;
  team: TeamData;
  remaining_time: number;
}
