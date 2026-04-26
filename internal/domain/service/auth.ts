'use client';

import { useState, useCallback } from 'react';
import { authRepository } from '../../adapters/outbound/repository/auth';
import { tokenStore } from '../../utils/token';

export function useAuth() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const signIn = useCallback(async (username: string, password: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      const tokens = await authRepository.signIn({ username, password });
      tokenStore.set(tokens.access_token, tokens.refresh_token);
      return true;
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Sign in failed');
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const signOut = useCallback(() => {
    tokenStore.clear();
    window.location.href = '/sign-in';
  }, []);

  const refresh = useCallback(async (): Promise<boolean> => {
    const refreshToken = tokenStore.getRefresh();
    if (!refreshToken) return false;
    try {
      const tokens = await authRepository.refresh(refreshToken);
      tokenStore.set(tokens.access_token, tokens.refresh_token);
      return true;
    } catch {
      tokenStore.clear();
      return false;
    }
  }, []);

  return { signIn, signOut, refresh, loading, error };
}
