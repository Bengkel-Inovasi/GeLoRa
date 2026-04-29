'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { alertRepository } from '../../../adapters/outbound/repository/alert';
import type { Alert } from '../../model/alert';

export function useAlerts(pollIntervalMs = 5000) {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const fetch = useCallback(async () => {
    try {
      const res = await alertRepository.list({ unacknowledged: true });
      setAlerts(res.data ?? []);
    } catch {
      // silently ignore — the listener should not crash the page
    }
  }, []);

  useEffect(() => {
    fetch();
    intervalRef.current = setInterval(fetch, pollIntervalMs);
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [fetch, pollIntervalMs]);

  const acknowledge = useCallback(async (id: number) => {
    await alertRepository.acknowledge(id);
    setAlerts((prev) => prev.filter((a) => a.id !== id));
  }, []);

  return { alerts, acknowledge, refetch: fetch };
}

export async function sendAlert(nodeId: number | null, message: string): Promise<void> {
  await alertRepository.create(nodeId, message);
}
