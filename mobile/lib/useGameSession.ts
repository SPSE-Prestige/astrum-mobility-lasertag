import { useEffect, useRef, useState, useCallback } from 'react';
import { getPlayerSession } from './api';
import type { PlayerSession } from './types';

interface UseGameSessionResult {
  session: PlayerSession | null;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
}

export function useGameSession(code: string): UseGameSessionResult {
  const [session, setSession] = useState<PlayerSession | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const mountedRef = useRef(true);

  const fetchSession = useCallback(async () => {
    try {
      const data = await getPlayerSession(code);
      if (!mountedRef.current) return;
      setSession(data);
      setError(null);

      if (data.game.status === 'finished' && intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    } catch (e) {
      if (!mountedRef.current) return;
      setError(e instanceof Error ? e.message : 'Chyba připojení');
    } finally {
      if (mountedRef.current) setLoading(false);
    }
  }, [code]);

  useEffect(() => {
    mountedRef.current = true;
    setLoading(true);
    setError(null);

    fetchSession();
    intervalRef.current = setInterval(fetchSession, 2000);

    return () => {
      mountedRef.current = false;
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [fetchSession]);

  return { session, loading, error, refresh: fetchSession };
}
