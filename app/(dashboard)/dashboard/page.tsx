'use client';

import dynamic from 'next/dynamic';
import { useState } from 'react';
import { useNodes } from '@/internal/domain/service/http/node';
import { useNodeRecords } from '@/internal/domain/service/http/record';
import NodeList from '@/components/dashboard/NodeList';
import NodePanel from '@/components/dashboard/NodePanel';
import Header from '@/components/layout/Header';
import Spinner from '@/components/ui/Spinner';
import { POLL_INTERVAL_MS } from '@/internal/config/constants';
import type { ComponentProps } from 'react';
import type MapViewType from '@/components/map/MapView';

const MapView = dynamic<ComponentProps<typeof MapViewType>>(
  () => import('@/components/map/MapView'),
  { ssr: false },
);

export default function DashboardPage() {
  const { nodes, loading } = useNodes();
  const nodeIds = nodes.map((n) => n.id);
  const { recordMap } = useNodeRecords(nodeIds);
  const [selectedNodeId, setSelectedNodeId] = useState<number | null>(null);

  const selectedNode = nodes.find((n) => n.id === selectedNodeId) ?? null;
  const selectedRecord = selectedNodeId != null ? recordMap.get(selectedNodeId) : undefined;

  return (
    <>
      <Header title="Live Tracker" />
      <div className="flex flex-1 overflow-hidden">
        <div className="w-64 border-r border-slate-200 bg-white overflow-hidden flex flex-col shrink-0">
          <NodeList
            nodes={nodes}
            recordMap={recordMap}
            selectedNodeId={selectedNodeId}
            onSelectNode={setSelectedNodeId}
            loading={loading}
          />
        </div>

        <div className="flex-1 relative">
          {loading && nodes.length === 0 ? (
            <div className="flex h-full items-center justify-center">
              <Spinner size={32} />
            </div>
          ) : (
            <MapView
              nodes={nodes}
              recordMap={recordMap}
              selectedNodeId={selectedNodeId}
              onSelectNode={setSelectedNodeId}
            />
          )}
          <div className="absolute bottom-4 right-4 z-[1000] bg-white/80 backdrop-blur text-xs text-slate-500 px-3 py-1 rounded-full shadow">
            Auto-refresh every {POLL_INTERVAL_MS / 1000}s
          </div>
        </div>

        <div className="w-72 border-l border-slate-200 bg-white p-4 flex flex-col shrink-0">
          <NodePanel node={selectedNode} record={selectedRecord} />
        </div>
      </div>
    </>
  );
}
