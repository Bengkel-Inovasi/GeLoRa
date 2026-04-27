export const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? 'http://localhost:8080';

export const TILE_URL =
  process.env.NEXT_PUBLIC_TILE_URL ??
  'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';

export const TILE_ATTRIBUTION =
  process.env.NEXT_PUBLIC_TILE_ATTRIBUTION ??
  '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors';

export const POLL_INTERVAL_MS = Number(process.env.NEXT_PUBLIC_POLL_INTERVAL_MS ?? 5000);

export const MAP_DEFAULT_CENTER: [number, number] = [-2.5, 118];
export const MAP_DEFAULT_ZOOM = 5;

export const RECORD_LOOKBACK_SECONDS = 300;

// How far back to fetch trail history (2 hours)
export const TRAIL_LOOKBACK_SECONDS = 7200;
// Trail data is fetched less often than live vitals
export const TRAIL_POLL_INTERVAL_MS = 30_000;
