export function formatHeartRate(bpm: number | null | undefined): string {
  if (bpm == null) return '—';
  return `${bpm.toFixed(0)} bpm`;
}

export function formatBodyTemp(celsius: number | null | undefined): string {
  if (celsius == null) return '—';
  return `${celsius.toFixed(1)} °C`;
}

export function formatCoords(lat: number | null | undefined, lon: number | null | undefined): string {
  if (lat == null || lon == null) return '—';
  const latStr = `${Math.abs(lat).toFixed(5)}° ${lat >= 0 ? 'N' : 'S'}`;
  const lonStr = `${Math.abs(lon).toFixed(5)}° ${lon >= 0 ? 'E' : 'W'}`;
  return `${latStr}, ${lonStr}`;
}

export function formatTime(iso: string | null | undefined): string {
  if (!iso) return '—';
  // Fixed format (HH:MM:SS WIB) — avoids server/client hydration mismatch
  // caused by toLocaleTimeString() returning different values per timezone
  const d = new Date(iso);
  const wib = new Date(d.getTime() + 7 * 60 * 60 * 1000);
  const hh = wib.getUTCHours().toString().padStart(2, '0');
  const mm = wib.getUTCMinutes().toString().padStart(2, '0');
  const ss = wib.getUTCSeconds().toString().padStart(2, '0');
  return `${hh}:${mm}:${ss} WIB`;
}

export function secondsAgo(seconds: number): string {
  const d = new Date(Date.now() - seconds * 1000);
  return d.toISOString();
}
