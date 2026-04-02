// ── API Error class with typed status codes ──

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    public readonly data?: unknown,
  ) {
    super(code || `HTTP ${status}`);
    this.name = "ApiError";
  }

  get isAuthError(): boolean {
    return this.status === 401;
  }

  get isNotFound(): boolean {
    return this.status === 404;
  }

  get isConflict(): boolean {
    return this.status === 409;
  }

  get isServerError(): boolean {
    return this.status >= 500;
  }
}

// ── Base HTTP client with timeout, error typing, token management ──

const DEFAULT_TIMEOUT_MS = 10_000;

const API_BASE =
  typeof window !== "undefined"
    ? (process.env.NEXT_PUBLIC_API_URL ?? "")
    : (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080");

export type TokenExpiredCallback = () => void;

class HttpClient {
  private token: string | null = null;
  private onTokenExpired: TokenExpiredCallback | null = null;

  setToken(token: string | null) {
    this.token = token;
  }

  setOnTokenExpired(cb: TokenExpiredCallback | null) {
    this.onTokenExpired = cb;
  }

  async request<T>(
    method: string,
    path: string,
    body?: unknown,
    timeoutMs: number = DEFAULT_TIMEOUT_MS,
  ): Promise<T> {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), timeoutMs);

    try {
      const headers: HeadersInit = { "Content-Type": "application/json" };
      if (this.token) {
        headers["Authorization"] = `Bearer ${this.token}`;
      }

      const res = await fetch(`${API_BASE}/api${path}`, {
        method,
        headers,
        body: body !== undefined ? JSON.stringify(body) : undefined,
        signal: controller.signal,
      });

      if (!res.ok) {
        const data = await res.json().catch(() => ({ error: `HTTP ${res.status}` }));
        const error = new ApiError(res.status, data.error || data.code || `HTTP ${res.status}`, data);

        if (error.isAuthError && this.onTokenExpired) {
          this.onTokenExpired();
        }

        throw error;
      }

      if (res.status === 204) return undefined as T;
      return res.json() as Promise<T>;
    } catch (err) {
      if (err instanceof ApiError) throw err;
      if (err instanceof DOMException && err.name === "AbortError") {
        throw new ApiError(0, "Request timeout", { timeoutMs });
      }
      throw new ApiError(0, err instanceof Error ? err.message : "Network error");
    } finally {
      clearTimeout(timeout);
    }
  }

  get = <T>(path: string) => this.request<T>("GET", path);
  post = <T>(path: string, body?: unknown) => this.request<T>("POST", path, body);
  patch = <T>(path: string, body?: unknown) => this.request<T>("PATCH", path, body);
  del = <T>(path: string) => this.request<T>("DELETE", path);
}

export const httpClient = new HttpClient();
