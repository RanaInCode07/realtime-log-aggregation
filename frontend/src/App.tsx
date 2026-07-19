import { useEffect, useMemo, useState } from "react";
import { useLogWebSocket } from "./hooks/useLogWebSocket";
import type { LogLevel } from "./types/log";

const levels: LogLevel[] = ["INFO", "DEBUG", "WARN", "ERROR"];

function formatTime(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  const time = date.toLocaleTimeString([], { hour12: false });
  return `${time}.${String(date.getMilliseconds()).padStart(3, "0")}`;
}

function App() {
  const { logs, connectionState, clearLogs } = useLogWebSocket();
  const [search, setSearch] = useState("");
  const [selectedLevel, setSelectedLevel] = useState<LogLevel | "ALL">("ALL");
  const [selectedService, setSelectedService] = useState("ALL");
  const [paused, setPaused] = useState(false);
  const [frozenLogs, setFrozenLogs] = useState<typeof logs>([]);

  useEffect(() => {
    if (!paused) setFrozenLogs(logs);
  }, [logs, paused]);

  const displayedLogs = paused ? frozenLogs : logs;
  const services = useMemo(() => [...new Set(logs.map((log) => log.service))].sort(), [logs]);
  const filteredLogs = useMemo(() => {
    const needle = search.trim().toLowerCase();
    return displayedLogs.filter((log) =>
      (selectedLevel === "ALL" || log.level === selectedLevel) &&
      (selectedService === "ALL" || log.service === selectedService) &&
      (!needle || `${log.message} ${log.service} ${log.level}`.toLowerCase().includes(needle)),
    );
  }, [displayedLogs, search, selectedLevel, selectedService]);

  const errors = logs.filter((log) => log.level === "ERROR").length;
  const warnings = logs.filter((log) => log.level === "WARN").length;

  function togglePause() {
    if (!paused) setFrozenLogs(logs);
    setPaused((current) => !current);
  }

  return (
    <main className="app-shell">
      <div className="orb orb-one" /><div className="orb orb-two" />
      <header className="topbar">
        <div className="brand"><span className="brand-mark">P</span><span>PulseLog</span><span className="environment">LOCAL</span></div>
        <div className={`connection ${connectionState}`}><span className="status-dot" />{connectionState === "connected" ? "Live stream connected" : connectionState === "connecting" ? "Connecting…" : "Reconnecting…"}</div>
      </header>

      <section className="hero">
        <div><p className="eyebrow">REAL-TIME OBSERVABILITY</p><h1>Logs, without the noise.</h1><p className="subheading">Monitor your event pipeline as it happens. Filter signal from noise and stay ahead of issues.</p></div>
        <button className="primary-button" onClick={togglePause}>{paused ? "Resume live feed" : "Pause live feed"}</button>
      </section>

      <section className="metrics" aria-label="Log metrics">
        <Metric label="Events received" value={logs.length.toLocaleString()} accent="blue" icon="↗" />
        <Metric label="Errors detected" value={errors.toLocaleString()} accent="red" icon="!" />
        <Metric label="Warnings" value={warnings.toLocaleString()} accent="amber" icon="△" />
        <Metric label="Active services" value={services.length.toLocaleString()} accent="violet" icon="◌" />
      </section>

      <section className="log-panel">
        <div className="panel-heading"><div><h2>Live event stream</h2><p>{filteredLogs.length} matching events {paused && "· Feed paused"}</p></div><button className="text-button" onClick={clearLogs}>Clear events</button></div>
        <div className="filters">
          <label className="search"><span>⌕</span><input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search messages, services, or levels" /></label>
          <select value={selectedLevel} onChange={(event) => setSelectedLevel(event.target.value as LogLevel | "ALL")}><option value="ALL">All levels</option>{levels.map((level) => <option key={level}>{level}</option>)}</select>
          <select value={selectedService} onChange={(event) => setSelectedService(event.target.value)}><option value="ALL">All services</option>{services.map((service) => <option key={service}>{service}</option>)}</select>
        </div>
        <div className="table-wrap"><table><thead><tr><th>Timestamp</th><th>Level</th><th>Service</th><th>Message</th></tr></thead><tbody>{filteredLogs.length ? filteredLogs.map((log) => <tr key={log.id}><td className="timestamp">{formatTime(log.time)}</td><td><span className={`badge ${log.level.toLowerCase()}`}>{log.level}</span></td><td><span className="service">{log.service}</span></td><td className="message">{log.message}</td></tr>) : <tr><td colSpan={4} className="empty"><span>No matching events yet.</span><small>Start the generator and make sure the consumer is running.</small></td></tr>}</tbody></table></div>
      </section>
    </main>
  );
}

function Metric({ label, value, accent, icon }: { label: string; value: string; accent: string; icon: string }) {
  return <article className={`metric-card ${accent}`}><div><p>{label}</p><strong>{value}</strong></div><span className="metric-icon">{icon}</span></article>;
}

export default App;
