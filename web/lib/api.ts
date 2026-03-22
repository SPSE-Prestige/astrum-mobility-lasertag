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
  id: string;
  game_id: string;
  name: string;
  color: string;
}

export interface PlayerResponse {
  id: string;
  game_id: string;
  team_id?: string;
  device_id: string;
  nickname: string;
  score: number;
  kills: number;
  deaths: number;
  is_alive: boolean;
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

// ── API Client ──

const API_BASE =
  typeof window !== "undefined"
    ? (process.env.NEXT_PUBLIC_API_URL ?? "")
    : (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080");

class ApiClient {
  private token: string | null = null;

  setToken(token: string | null) {
    this.token = token;
  }

  private async request<T>(
    method: string,
    path: string,
    body?: unknown,
  ): Promise<T> {
    const headers: HeadersInit = { "Content-Type": "application/json" };
    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const res = await fetch(`${API_BASE}/api${path}`, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      const data = await res
        .json()
        .catch(() => ({ error: `HTTP ${res.status}` }));
      throw new Error(data.error || `HTTP ${res.status}`);
    }

    if (res.status === 204) return undefined as T;
    return res.json() as Promise<T>;
  }

  // ── Auth ──

  async login(
    username: string,
    password: string,
  ): Promise<LoginResponse> {
    return this.request("POST", "/auth/login", { username, password });
  }

  async logout(): Promise<void> {
    return this.request("POST", "/auth/logout");
  }

  // ── Devices ──

  async listDevices(): Promise<DeviceResponse[]> {
    return this.request("GET", "/devices");
  }

  async listAvailableDevices(): Promise<DeviceResponse[]> {
    return this.request("GET", "/devices/available");
  }

  // ── Games ──

  async createGame(settings?: GameSettingsDTO): Promise<GameResponse> {
    return this.request("POST", "/games", { settings });
  }

  async getGame(id: string): Promise<GameResponse> {
    return this.request("GET", `/games/${encodeURIComponent(id)}`);
  }

  async getGameFull(id: string): Promise<GameFullResponse> {
    return this.request(
      "GET",
      `/games/${encodeURIComponent(id)}/full`,
    );
  }

  async listGames(): Promise<GameResponse[]> {
    return this.request("GET", "/games");
  }

  async startGame(id: string): Promise<GameResponse> {
    return this.request(
      "POST",
      `/games/${encodeURIComponent(id)}/start`,
    );
  }

  async endGame(id: string): Promise<GameResponse> {
    return this.request(
      "POST",
      `/games/${encodeURIComponent(id)}/end`,
    );
  }

  async updateSettings(
    id: string,
    settings: GameSettingsDTO,
  ): Promise<GameResponse> {
    return this.request(
      "PATCH",
      `/games/${encodeURIComponent(id)}/settings`,
      { settings },
    );
  }

  // ── Teams ──

  async addTeam(
    gameId: string,
    name: string,
    color: string,
  ): Promise<TeamResponse> {
    return this.request(
      "POST",
      `/games/${encodeURIComponent(gameId)}/teams`,
      { name, color },
    );
  }

  async listTeams(gameId: string): Promise<TeamResponse[]> {
    return this.request(
      "GET",
      `/games/${encodeURIComponent(gameId)}/teams`,
    );
  }

  async removeTeam(gameId: string, teamId: string): Promise<void> {
    return this.request(
      "DELETE",
      `/games/${encodeURIComponent(gameId)}/teams/${encodeURIComponent(teamId)}`,
    );
  }

  // ── Players ──

  async addPlayer(
    gameId: string,
    deviceId: string,
    nickname: string,
    teamId?: string,
  ): Promise<PlayerResponse> {
    return this.request(
      "POST",
      `/games/${encodeURIComponent(gameId)}/players`,
      { device_id: deviceId, nickname, team_id: teamId ?? null },
    );
  }

  async listPlayers(gameId: string): Promise<PlayerResponse[]> {
    return this.request(
      "GET",
      `/games/${encodeURIComponent(gameId)}/players`,
    );
  }

  async removePlayer(
    gameId: string,
    playerId: string,
  ): Promise<void> {
    return this.request(
      "DELETE",
      `/games/${encodeURIComponent(gameId)}/players/${encodeURIComponent(playerId)}`,
    );
  }

  async updatePlayerTeam(
    gameId: string,
    playerId: string,
    teamId: string | null,
  ): Promise<void> {
    return this.request(
      "PATCH",
      `/games/${encodeURIComponent(gameId)}/players/${encodeURIComponent(playerId)}/team`,
      { team_id: teamId },
    );
  }

  // ── Leaderboard & Events ──

  async getLeaderboard(gameId: string): Promise<PlayerResponse[]> {
    return this.request(
      "GET",
      `/games/${encodeURIComponent(gameId)}/leaderboard`,
    );
  }

  async getEvents(gameId: string): Promise<EventResponse[]> {
    return this.request(
      "GET",
      `/games/${encodeURIComponent(gameId)}/events`,
    );
  }
}

export const api = new ApiClient();
