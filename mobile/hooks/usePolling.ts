import { useCallback, useEffect, useRef } from "react";
import { POLLING_INTERVAL_MS } from "@/constants/config";

interface UsePollingOptions {
  enabled: boolean;
  interval?: number;
  onTick: () => Promise<void> | void;
}

export function usePolling({ enabled, interval = POLLING_INTERVAL_MS, onTick }: UsePollingOptions) {
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const tickRef = useRef(onTick);
  tickRef.current = onTick;

  const stop = useCallback(() => {
    if (timerRef.current !== null) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  useEffect(() => {
    if (!enabled) {
      stop();
      return;
    }

    timerRef.current = setInterval(() => {
      tickRef.current();
    }, interval);

    return stop;
  }, [enabled, interval, stop]);
}

export default usePolling;
