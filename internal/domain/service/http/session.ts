'use client';

import { useState, useEffect, useCallback } from 'react';
import { sessionRepository } from '../../../adapters/outbound/repository/session';
import type { Session } from '../../model/session';

export function useSessions() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchSessions = useCallback(async () => {
    try {
      const res = await sessionRepository.list({ limit: 200 });
      setSessions(res.data);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load sessions');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSessions();
  }, [fetchSessions]);

  const activeSessions = sessions.filter((s) => s.ended_at == null);

  return { sessions, activeSessions, loading, error, refetch: fetchSessions };
}
