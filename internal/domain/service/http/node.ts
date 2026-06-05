'use client';

import { useState, useEffect, useCallback } from 'react';
import { nodeRepository } from '../../../adapters/outbound/repository/node';
import type { NodeListItem, Node } from '../../model/node';

export function useNodes() {
  const [nodes, setNodes] = useState<NodeListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchNodes = useCallback(async () => {
    try {
      const res = await nodeRepository.list({ limit: 100 });
      setNodes(res.data);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load nodes');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchNodes();
    // Poll every 10 seconds to keep aliases and device list in sync across components
    const id = setInterval(fetchNodes, 10000);
    return () => clearInterval(id);
  }, [fetchNodes]);

  return { nodes, loading, error, refetch: fetchNodes };
}

export function useNode(id: number | null) {
  const [node, setNode] = useState<Node | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (id == null) { setNode(null); return; }
    setLoading(true);
    nodeRepository.getById(id)
      .then(setNode)
      .catch(() => setNode(null))
      .finally(() => setLoading(false));
  }, [id]);

  return { node, loading };
}
