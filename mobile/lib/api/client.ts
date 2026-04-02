import { API_BASE_URL, REQUEST_TIMEOUT_MS } from "@/constants/config";
import type { PlayerSessionResponse } from "./types";

class ApiClient {
  private baseUrl: string;
  private timeout: number;

  constructor(baseUrl: string, timeout: number) {
    this.baseUrl = baseUrl.replace(/\/+$/, "");
    this.timeout = timeout;
  }

  private async request<T>(path: string): Promise<T> {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeout);

    try {
      const url = `${this.baseUrl}${path}`;
      const response = await fetch(url, {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        signal: controller.signal,
      });

      if (!response.ok) {
        const body = await response.json().catch(() => null);
        const code = body?.error?.code ?? "UNKNOWN";
        const message = body?.error?.message ?? response.statusText;
        throw new ApiRequestError(response.status, code, message);
      }

      return (await response.json()) as T;
    } catch (error) {
      if (error instanceof ApiRequestError) throw error;
      if (error instanceof DOMException && error.name === "AbortError") {
        throw new ApiRequestError(0, "TIMEOUT", "Request timed out");
      }
      throw new ApiRequestError(0, "NETWORK_ERROR", "Unable to connect to server");
    } finally {
      clearTimeout(timer);
    }
  }

  async getPlayerSession(code: string): Promise<PlayerSessionResponse> {
    return this.request<PlayerSessionResponse>(
      `/api/player/session/${encodeURIComponent(code.toUpperCase())}`
    );
  }

  async healthCheck(): Promise<boolean> {
    try {
      await this.request<unknown>("/health");
      return true;
    } catch {
      return false;
    }
  }
}

export class ApiRequestError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string
  ) {
    super(message);
    this.name = "ApiRequestError";
  }

  get isNotFound(): boolean {
    return this.status === 404;
  }

  get isNetwork(): boolean {
    return this.status === 0;
  }
}

export const apiClient = new ApiClient(API_BASE_URL, REQUEST_TIMEOUT_MS);
export default apiClient;
