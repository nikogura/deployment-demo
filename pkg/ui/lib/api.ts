import type { AppInfo } from "@/types";

// eslint-disable-next-line local-rules/disallow-fetch
export async function fetchInfo(): Promise<AppInfo> {
  const resp = await fetch("/api/info");
  const data = (await resp.json()) as AppInfo;
  return data;
}
