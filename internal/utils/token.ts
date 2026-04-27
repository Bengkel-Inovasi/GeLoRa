const ACCESS_KEY = 'access_token';
const REFRESH_KEY = 'refresh_token';

function decodeJwtPayload(token: string): Record<string, unknown> | null {
  try {
    const b64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(b64));
  } catch {
    return null;
  }
}

export const tokenStore = {
  getAccess: (): string | null => localStorage.getItem(ACCESS_KEY),
  getRefresh: (): string | null => localStorage.getItem(REFRESH_KEY),
  set: (access: string, refresh: string): void => {
    localStorage.setItem(ACCESS_KEY, access);
    localStorage.setItem(REFRESH_KEY, refresh);
  },
  clear: (): void => {
    localStorage.removeItem(ACCESS_KEY);
    localStorage.removeItem(REFRESH_KEY);
  },
  hasToken: (): boolean => !!localStorage.getItem(ACCESS_KEY),
  // Decode JWT payload to read role without an API round-trip
  // Go's golang-jwt encodes struct field names as-is (no json tags → "Role", "Id")
  getRole: (): 'super' | 'admin' | 'client' | null => {
    const token = localStorage.getItem(ACCESS_KEY);
    if (!token) return null;
    const payload = decodeJwtPayload(token);
    const role = payload?.['Role'];
    if (role === 'super' || role === 'admin' || role === 'client') return role;
    return null;
  },
  getUserId: (): number | null => {
    const token = localStorage.getItem(ACCESS_KEY);
    if (!token) return null;
    const payload = decodeJwtPayload(token);
    const id = payload?.['Id'];
    return typeof id === 'number' ? id : null;
  },
};
