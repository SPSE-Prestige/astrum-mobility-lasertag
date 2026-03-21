"use client";

import { useEffect, useRef } from "react";

/**
 * Connects to the backend WebSocket during live phase.
 * Calls `onEvent` when a message matching the current gameId arrives,
 * which the caller uses to trigger an immediate poll.
 */
export function useWebSocket(
  enabled: boolean,
  gameId: string | null,
  onEvent: () => void,
) {
  const onEventRef = useRef(onEvent);
  onEventRef.current = onEvent;

  useEffect(() => {
    if (!enabled || !gameId) return;

    const apiUrl =
      process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const wsUrl = apiUrl.replace(/^http/, "ws") + "/ws";

    let ws: WebSocket;
    let reconnectTimer: ReturnType<typeof setTimeout>;
    let disposed = false;

    const connect = () => {
      if (disposed) return;
      ws = new WebSocket(wsUrl);

      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data);
          // Filter by game_id if present in the message
          if (msg.game_id && msg.game_id !== gameId) return;
          onEventRef.current();
        } catch {
          // ignore malformed messages
        }
      };

      ws.onclose = () => {
        if (!disposed) {
          reconnectTimer = setTimeout(connect, 2000);
        }
      };

      ws.onerror = () => {
        ws.close();
      };
    };

    connect();

    return () => {
      disposed = true;
      clearTimeout(reconnectTimer);
      ws?.close();
    };
  }, [enabled, gameId]);
}
