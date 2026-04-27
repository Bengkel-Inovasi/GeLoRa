'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { recordRepository } from '../../../adapters/outbound/repository/record';
import type { Record } from '../../model/record';
import { POLL_INTERVAL_MS, RECORD_LOOKBACK_SECONDS, TRAIL_LOOKBACK_SECONDS, TRAIL_POLL_INTERVAL_MS } from '../../../config/constants';
import { secondsAgo } from '../../../utils/format';

export type NodeRecordMap = Map<number, Record>;

// [lat, lng] pairs for a single node's trail, ordered oldest → newest
export type NodeTrailMap = Map<number, [number, number][]>;

function latestRecord(records: Record[]): Record | undefined {
  if (records.length === 0) return undefined;
  return records.reduce((a, b) => (new Date(a.time) > new Date(b.time) ? a : b));
}

export function useNodeRecords(nodeIds: number[]) {
  const [recordMap, setRecordMap] = useState<NodeRecordMap>(new Map());
  const nodeIdsRef = useRef(nodeIds);
  nodeIdsRef.current = nodeIds;

  const poll = useCallback(async () => {
    if (nodeIdsRef.current.length === 0) return;
    const start = secondsAgo(RECORD_LOOKBACK_SECONDS);
    const results = await Promise.allSettled(
      nodeIdsRef.current.map((id) =>
        recordRepository.list({ node_id: id, start_time: start }),
      ),
    );
    setRecordMap((prev) => {
      const next = new Map(prev);
      results.forEach((r, i) => {
        if (r.status === 'fulfilled') {
          const rec = latestRecord(r.value.data);
          if (rec) next.set(nodeIdsRef.current[i], rec);
        }
      });
      return next;
    });
  }, []);

  useEffect(() => {
    poll();
    const id = setInterval(poll, POLL_INTERVAL_MS);
    return () => clearInterval(id);
  }, [poll]);

  return { recordMap };
}

export function useNodeTrails(nodeIds: number[]) {
  const [trailMap, setTrailMap] = useState<NodeTrailMap>(new Map());
  const nodeIdsRef = useRef(nodeIds);
  nodeIdsRef.current = nodeIds;

  const fetch = useCallback(async () => {
    if (nodeIdsRef.current.length === 0) return;
    const start = secondsAgo(TRAIL_LOOKBACK_SECONDS);
    const results = await Promise.allSettled(
      nodeIdsRef.current.map((id) =>
        recordRepository.list({ node_id: id, start_time: start }),
      ),
    );
    setTrailMap(() => {
      const next = new Map<number, [number, number][]>();
      results.forEach((r, i) => {
        if (r.status !== 'fulfilled') return;
        const points: [number, number][] = r.value.data
          .filter((rec) => rec.latitude != null && rec.longitude != null)
          .map((rec) => [rec.latitude as number, rec.longitude as number]);
        if (points.length > 0) next.set(nodeIdsRef.current[i], points);
      });
      return next;
    });
  }, []);

  useEffect(() => {
    fetch();
    const id = setInterval(fetch, TRAIL_POLL_INTERVAL_MS);
    return () => clearInterval(id);
  }, [fetch]);

  return { trailMap };
}
