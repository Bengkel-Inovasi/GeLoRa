'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { recordRepository } from '../../../adapters/outbound/repository/record';
import type { Record } from '../../model/record';
import { POLL_INTERVAL_MS, RECORD_LOOKBACK_SECONDS } from '../../../config/constants';
import { secondsAgo } from '../../../utils/format';

export type NodeRecordMap = Map<number, Record>;

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
