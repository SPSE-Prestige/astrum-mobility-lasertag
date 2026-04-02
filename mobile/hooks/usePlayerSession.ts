import { useCallback, useEffect, useRef, useState } from "react";
import apiClient, { ApiRequestError } from "@/lib/api/client";
import type { PlayerSessionResponse } from "@/lib/api/types";

interface UsePlayerSessionReturn {
  data: PlayerSessionResponse | null;
  loading: boolean;
  error: string | null;
  errorCode: string | null;
  fetch: (code: string) => Promise<void>;
  refresh: () => Promise<void>;
  clear: () => void;
}

export function usePlayerSession(): UsePlayerSessionReturn {
  const [data, setData] = useState<PlayerSessionResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [errorCode, setErrorCode] = useState<string | null>(null);
  const lastCode = useRef<string | null>(null);

  const fetchSession = useCallback(async (code: string) => {
    lastCode.current = code;
    setLoading(true);
    setError(null);
    setErrorCode(null);

    try {
      const session = await apiClient.getPlayerSession(code);
      setData(session);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
        setErrorCode(err.code);
      } else {
        setError("An unexpected error occurred");
        setErrorCode("UNKNOWN");
      }
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const refresh = useCallback(async () => {
    if (lastCode.current) {
      await fetchSession(lastCode.current);
    }
  }, [fetchSession]);

  const clear = useCallback(() => {
    setData(null);
    setError(null);
    setErrorCode(null);
    lastCode.current = null;
  }, []);

  return { data, loading, error, errorCode, fetch: fetchSession, refresh, clear };
}

export default usePlayerSession;
