"use client";

import { useEffect, useRef } from "react";

interface UseWebSocketOptions {
  enabled: boolean;
  gameId: string | null;
  onEvent: () => void;
  onStatusChange?: (status: "connecting" | "connected" | "disconnected") => void;
}

const MAX_RECONNECT_DELAY = 30_000;
const INITIAL_RECONNECT_DELAY = 1_000;

/**
 * Connects to backend WebSocket during live phase.
 * Features: exponential backoff reconnect, connection status, typed events.
 */
export function useWebSocket({ enabled, gameId, onEvent, onStatusChange }: UseWebSocketOptions) {
  const onEventRef = useRef(onEvent);
  onEventRef.current = onEvent;

  const onStatusRef = useRef(onStatusChange);
  onStatusRef.current = onStatusChange;

  useEffect(() => {
    if (!enabled || !gameId) {
      onStatusRef.current?.("disconnected");
      return;
    }

    // Always use current page origin for WS (goes through nginx proxy)
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${proto}//${window.location.host}/ws`;

    let ws: WebSocket;
    let reconnectTimer: ReturnType<typeof setTimeout>;
    let disposed = false;
    let reconnectDelay = INITIAL_RECONNECT_DELAY;

    const connect = () => {
      if (disposed) return;
      onStatusRef.current?.("connecting");

      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        reconnectDelay = INITIAL_RECONNECT_DELAY;
        onStatusRef.current?.("connected");
      };

      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data);
          if (msg.game_id && msg.game_id !== gameId) return;
          onEventRef.current();
        } catch {
          // ignore malformed messages
        }
      };

      ws.onclose = () => {
        onStatusRef.current?.("disconnected");
        if (!disposed) {
          reconnectTimer = setTimeout(connect, reconnectDelay);
          reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY);
        }
      };

      ws.onerror = (event) => {
        console.warn("[WebSocket] Connection error", event);
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
