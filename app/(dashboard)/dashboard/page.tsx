'use client';

import dynamic from 'next/dynamic';
import { useState } from 'react';
import { useNodes } from '@/internal/domain/service/http/node';
import { useNodeRecords, useNodeTrails } from '@/internal/domain/service/http/record';
import NodeList from '@/components/dashboard/NodeList';
import NodePanel from '@/components/dashboard/NodePanel';
import Header from '@/components/layout/Header';
import Badge from '@/components/ui/Badge';
import Spinner from '@/components/ui/Spinner';
import { POLL_INTERVAL_MS } from '@/internal/config/constants';
import { formatHeartRate, formatBodyTemp, formatCoords, formatTime } from '@/internal/utils/format';
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
  const { trailMap } = useNodeTrails(nodeIds);
  const [selectedNodeId, setSelectedNodeId] = useState<number | null>(null);
  const [mobileTab, setMobileTab] = useState<'map' | 'list'>('map');

  const selectedNode = nodes.find((n) => n.id === selectedNodeId) ?? null;
  const selectedRecord = selectedNodeId != null ? recordMap.get(selectedNodeId) : undefined;

  function handleSelectNode(id: number) {
    setSelectedNodeId(id);
    setMobileTab('map');
  }

  return (
    <>
      <Header title="Live Tracker" />

      {/* ── Desktop layout (md+): 3-column ── */}
      <div className="hidden md:flex flex-1 overflow-hidden">
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
              trailMap={trailMap}
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

      {/* ── Mobile layout (<md): tab-based ── */}
      <div className="flex md:hidden flex-col flex-1 overflow-hidden">
        {/* Tab switcher */}
        <div className="flex bg-white border-b border-slate-200 shrink-0">
          <button
            onClick={() => setMobileTab('map')}
            className={`flex-1 py-2.5 text-sm font-medium transition-colors ${
              mobileTab === 'map'
                ? 'text-emerald-600 border-b-2 border-emerald-500'
                : 'text-slate-500'
            }`}
          >
            🗺️ Map
          </button>
          <button
            onClick={() => setMobileTab('list')}
            className={`flex-1 py-2.5 text-sm font-medium transition-colors ${
              mobileTab === 'list'
                ? 'text-emerald-600 border-b-2 border-emerald-500'
                : 'text-slate-500'
            }`}
          >
            🧗 Climbers
          </button>
        </div>

        {mobileTab === 'map' ? (
          <div className="flex-1 flex flex-col overflow-hidden">
            {/* Map fills remaining space */}
            <div className="flex-1 relative min-h-0">
              {loading && nodes.length === 0 ? (
                <div className="flex h-full items-center justify-center">
                  <Spinner size={32} />
                </div>
              ) : (
                <MapView
                  nodes={nodes}
                  recordMap={recordMap}
                  trailMap={trailMap}
                  selectedNodeId={selectedNodeId}
                  onSelectNode={setSelectedNodeId}
                />
              )}
              <div className="absolute bottom-2 right-2 z-[1000] bg-white/80 backdrop-blur text-xs text-slate-500 px-2 py-1 rounded-full shadow">
                Refresh {POLL_INTERVAL_MS / 1000}s
              </div>
            </div>

            {/* Selected node panel — slides up from bottom when a node is picked */}
            {selectedNode ? (
              <div className="bg-white border-t border-slate-200 p-3 shrink-0">
                <div className="flex items-start justify-between mb-2">
                  <div className="min-w-0 flex-1">
                    <p className="font-semibold text-slate-800 text-sm truncate">{selectedNode.name}</p>
                    <p className="text-xs text-slate-400 font-mono">{selectedNode.mid}</p>
                  </div>
                  <div className="flex items-center gap-2 ml-2 shrink-0">
                    <Badge
                      label={selectedNode.is_validated ? 'Validated' : 'Unvalidated'}
                      variant={selectedNode.is_validated ? 'green' : 'red'}
                    />
                    <button
                      onClick={() => setSelectedNodeId(null)}
                      className="text-slate-400 hover:text-slate-600 text-xl leading-none"
                      aria-label="Close"
                    >
                      ×
                    </button>
                  </div>
                </div>

                {selectedRecord ? (
                  <div className="grid grid-cols-2 gap-2">
                    <div className="bg-slate-50 rounded-lg px-3 py-2">
                      <p className="text-xs text-slate-400">Heart Rate</p>
                      <p className="text-sm font-semibold text-slate-800">{formatHeartRate(selectedRecord.heart_rate)}</p>
                    </div>
                    <div className="bg-slate-50 rounded-lg px-3 py-2">
                      <p className="text-xs text-slate-400">Body Temp</p>
                      <p className="text-sm font-semibold text-slate-800">{formatBodyTemp(selectedRecord.body_temperature)}</p>
                    </div>
                    <div className="bg-slate-50 rounded-lg px-3 py-2 col-span-2">
                      <p className="text-xs text-slate-400">Coordinates · {formatTime(selectedRecord.time)}</p>
                      <p className="text-sm font-semibold text-slate-800 font-mono">{formatCoords(selectedRecord.latitude, selectedRecord.longitude)}</p>
                    </div>
                  </div>
                ) : (
                  <p className="text-sm text-slate-400 italic">No recent data received</p>
                )}
              </div>
            ) : (
              <div className="bg-white border-t border-slate-200 px-4 py-3 shrink-0">
                <p className="text-xs text-slate-400 text-center">Tap a marker on the map to view climber data</p>
              </div>
            )}
          </div>
        ) : (
          /* List tab */
          <div className="flex-1 overflow-hidden bg-white">
            <NodeList
              nodes={nodes}
              recordMap={recordMap}
              selectedNodeId={selectedNodeId}
              onSelectNode={handleSelectNode}
              loading={loading}
            />
          </div>
        )}
      </div>
    </>
  );
}
