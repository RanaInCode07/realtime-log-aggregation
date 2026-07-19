import { useCallback, useEffect, useRef, useState } from "react";
import type { ConnectionState, LogEvent } from "../types/log";

const MAX_LOGS = 1_000;
const RECONNECT_DELAY = 2_000;

export function useLogWebSocket() {
  const [logs, setLogs] = useState<LogEvent[]>([]);
  const [connectionState, setConnectionState] = useState<ConnectionState>("connecting");
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<number | null>(null);
  const shouldReconnectRef = useRef(true);

  useEffect(() => {
    const url = import.meta.env.VITE_WS_URL ?? "ws://localhost:8082/ws";

    const connect = () => {
      setConnectionState("connecting");
      const socket = new WebSocket(url);
      socketRef.current = socket;

      socket.onopen = () => setConnectionState("connected");
      socket.onmessage = (event) => {
        try {
          const incoming = JSON.parse(event.data) as LogEvent;
          setLogs((current) => [incoming, ...current].slice(0, MAX_LOGS));
        } catch {
          // Ignore malformed messages and keep the live stream available.
        }
      };
      socket.onerror = () => socket.close();
      socket.onclose = () => {
        setConnectionState("disconnected");
        if (shouldReconnectRef.current) {
          reconnectTimerRef.current = window.setTimeout(connect, RECONNECT_DELAY);
        }
      };
    };

    connect();
    return () => {
      shouldReconnectRef.current = false;
      if (reconnectTimerRef.current) window.clearTimeout(reconnectTimerRef.current);
      socketRef.current?.close();
    };
  }, []);

  const clearLogs = useCallback(() => setLogs([]), []);
  return { logs, connectionState, clearLogs };
}
