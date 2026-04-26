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
  return new Date(iso).toLocaleTimeString();
}

export function secondsAgo(seconds: number): string {
  const d = new Date(Date.now() - seconds * 1000);
  return d.toISOString();
}
