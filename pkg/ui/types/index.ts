export interface AppInfo {
  readonly version: string;
  readonly theme: string;
  readonly health: string;
  readonly healthy: boolean;
  readonly buildTime: string;
  readonly uptime: string;
  readonly uptimeSec: number;
}
