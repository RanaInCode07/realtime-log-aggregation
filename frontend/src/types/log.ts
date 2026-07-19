export type LogLevel = "INFO" | "DEBUG" | "WARN" | "ERROR";

export interface LogEvent {
  id: string;
  time: string;
  level: LogLevel;
  service: string;
  message: string;
}

export type ConnectionState = "connecting" | "connected" | "disconnected";
