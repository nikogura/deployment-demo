"use client";

import { AlertTriangle, CheckCircle, RefreshCw } from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";

import { fetchInfo } from "@/lib/api";
import type { AppInfo } from "@/types";

const REFRESH_OPTIONS: readonly { readonly label: string; readonly ms: number }[] = [
  { label: "1s", ms: 1000 },
  { label: "5s", ms: 5000 },
  { label: "10s", ms: 10000 },
  { label: "30s", ms: 30000 },
];

export default function DemoPage(): React.ReactElement {
  const [info, setInfo] = useState<AppInfo | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [refreshMs, setRefreshMs] = useState(5000);
  const [lastFetch, setLastFetch] = useState<Date | null>(null);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const loadInfo = useCallback(async (): Promise<void> => {
    try {
      const data = await fetchInfo();
      setInfo(data);
      setError(null);
      setLastFetch(new Date());
      document.documentElement.setAttribute("data-theme", data.theme);
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Failed to fetch";
      setError(msg);
    }
  }, []);

  useEffect((): (() => void) => {
    void loadInfo();

    intervalRef.current = setInterval((): void => {
      void loadInfo();
    }, refreshMs);

    return (): void => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [loadInfo, refreshMs]);

  const handleRefreshChange = (e: React.ChangeEvent<HTMLSelectElement>): void => {
    setRefreshMs(Number(e.target.value));
  };

  return (
    <main style={{
      minHeight: "100vh",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      padding: "2rem",
      position: "relative",
    }}>
      {/* Broken banner */}
      {info !== null && !info.healthy && (
        <div style={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          padding: "0.75rem",
          backgroundColor: "var(--accent)",
          color: "#fff",
          textAlign: "center",
          fontWeight: 700,
          fontSize: "1.1rem",
          zIndex: 100,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          gap: "0.5rem",
        }}>
          <AlertTriangle size={20} />
          THIS VERSION IS BROKEN — /healthz returning 503 — EXPECTING ROLLBACK
          <AlertTriangle size={20} />
        </div>
      )}

      {/* Version badge — the hero element */}
      <div style={{
        fontSize: "8rem",
        fontWeight: 900,
        fontFamily: "monospace",
        color: "var(--accent)",
        textShadow: "0 0 80px var(--accent-glow), 0 0 160px var(--accent-glow)",
        lineHeight: 1,
        marginBottom: "1rem",
        userSelect: "none",
        letterSpacing: "-0.05em",
      }}>
        {info?.version ?? "..."}
      </div>

      {/* Health indicator */}
      <div style={{
        display: "flex",
        alignItems: "center",
        gap: "0.5rem",
        fontSize: "1.5rem",
        fontWeight: 600,
        color: info?.healthy ? "#22c55e" : "#ef4444",
        marginBottom: "2rem",
      }}>
        {info?.healthy ? <CheckCircle size={28} /> : <AlertTriangle size={28} />}
        {info?.healthy ? "HEALTHY" : "UNHEALTHY"}
      </div>

      {/* Info card */}
      <div style={{
        backgroundColor: "var(--bg-card)",
        border: "1px solid var(--border)",
        borderRadius: "1rem",
        padding: "2rem",
        width: "100%",
        maxWidth: "480px",
        boxShadow: "0 0 40px var(--accent-glow)",
      }}>
        <table style={{ width: "100%", borderCollapse: "collapse", fontSize: "0.95rem" }}>
          <tbody>
            {[
              ["Version", info?.version],
              ["Theme", info?.theme],
              ["Health", info?.health],
              ["Build Time", info?.buildTime],
              ["Uptime", info?.uptime],
            ].map(([label, value]) => (
              <tr key={String(label)} style={{ borderBottom: "1px solid var(--border)" }}>
                <td style={{ padding: "0.6rem 0", color: "var(--text-muted)", fontWeight: 500 }}>{label}</td>
                <td style={{ padding: "0.6rem 0", textAlign: "right", fontFamily: "monospace" }}>{value ?? "\u2014"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Refresh control */}
      <div style={{
        marginTop: "1.5rem",
        display: "flex",
        alignItems: "center",
        gap: "0.75rem",
        color: "var(--text-muted)",
        fontSize: "0.85rem",
      }}>
        <RefreshCw size={14} />
        <span>Refresh:</span>
        <select
          value={refreshMs}
          onChange={handleRefreshChange}
          style={{
            padding: "0.25rem 0.5rem",
            backgroundColor: "var(--bg-card)",
            border: "1px solid var(--border)",
            borderRadius: "0.25rem",
            color: "var(--text)",
            fontSize: "0.85rem",
          }}
        >
          {REFRESH_OPTIONS.map((opt) => (
            <option key={opt.ms} value={opt.ms}>{opt.label}</option>
          ))}
        </select>
        {lastFetch !== null && (
          <span>last: {lastFetch.toLocaleTimeString()}</span>
        )}
      </div>

      {/* Error display */}
      {error !== null && (
        <div style={{
          marginTop: "1rem",
          padding: "0.75rem 1rem",
          borderRadius: "0.5rem",
          backgroundColor: "rgba(239, 68, 68, 0.1)",
          border: "1px solid #ef4444",
          color: "#ef4444",
          fontSize: "0.85rem",
          maxWidth: "480px",
          width: "100%",
          textAlign: "center",
        }}>
          {error}
        </div>
      )}

      {/* Footer */}
      <div style={{
        position: "fixed",
        bottom: "1rem",
        color: "var(--text-muted)",
        fontSize: "0.75rem",
      }}>
        deployment-demo — observability + deployment lifecycle demonstration
      </div>
    </main>
  );
}
