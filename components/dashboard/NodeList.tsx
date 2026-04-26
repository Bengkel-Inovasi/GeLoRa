'use client';

import type { NodeListItem } from '../../internal/domain/model/node';
import type { NodeRecordMap } from '../../internal/domain/service/http/record';
import Badge from '../ui/Badge';
import Spinner from '../ui/Spinner';
import { formatHeartRate, formatBodyTemp, formatTime } from '../../internal/utils/format';

interface Props {
  nodes: NodeListItem[];
  recordMap: NodeRecordMap;
  selectedNodeId: number | null;
  onSelectNode: (id: number) => void;
  loading: boolean;
}

export default function NodeList({ nodes, recordMap, selectedNodeId, onSelectNode, loading }: Props) {
  const activeNodes = nodes.filter((n) => recordMap.has(n.id));
  const inactiveNodes = nodes.filter((n) => !recordMap.has(n.id));

  return (
    <div className="flex flex-col h-full">
      <div className="px-4 py-3 border-b border-slate-200">
        <h2 className="text-sm font-semibold text-slate-700">Devices</h2>
        <p className="text-xs text-slate-400 mt-0.5">{activeNodes.length} active · {nodes.length} total</p>
      </div>

      <div className="overflow-y-auto flex-1">
        {loading && (
          <div className="flex justify-center py-8">
            <Spinner />
          </div>
        )}

        {activeNodes.length > 0 && (
          <div>
            <p className="px-4 py-2 text-xs text-slate-400 uppercase tracking-wide font-medium bg-slate-50">Active</p>
            {activeNodes.map((node) => {
              const rec = recordMap.get(node.id);
              const isSelected = node.id === selectedNodeId;
              return (
                <button
                  key={node.id}
                  onClick={() => onSelectNode(node.id)}
                  className={`w-full text-left px-4 py-3 border-b border-slate-100 hover:bg-emerald-50 transition-colors ${isSelected ? 'bg-emerald-50 border-l-4 border-l-emerald-500' : ''}`}
                >
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-medium text-slate-800">{node.name}</p>
                    <Badge label="Active" variant="green" />
                  </div>
                  <p className="text-xs text-slate-400 font-mono mt-0.5">{node.mid}</p>
                  {rec && (
                    <div className="flex gap-3 mt-1.5 text-xs text-slate-500">
                      <span>HR: {formatHeartRate(rec.heart_rate)}</span>
                      <span>Temp: {formatBodyTemp(rec.body_temperature)}</span>
                    </div>
                  )}
                  {rec && <p className="text-xs text-slate-400 mt-1">{formatTime(rec.time)}</p>}
                </button>
              );
            })}
          </div>
        )}

        {inactiveNodes.length > 0 && (
          <div>
            <p className="px-4 py-2 text-xs text-slate-400 uppercase tracking-wide font-medium bg-slate-50">No Recent Data</p>
            {inactiveNodes.map((node) => (
              <button
                key={node.id}
                onClick={() => onSelectNode(node.id)}
                className={`w-full text-left px-4 py-3 border-b border-slate-100 hover:bg-slate-50 opacity-60 transition-colors ${node.id === selectedNodeId ? 'bg-slate-100' : ''}`}
              >
                <div className="flex items-center justify-between">
                  <p className="text-sm font-medium text-slate-600">{node.name}</p>
                  <Badge label={node.is_validated ? 'Validated' : 'Unvalidated'} variant={node.is_validated ? 'blue' : 'gray'} />
                </div>
                <p className="text-xs text-slate-400 font-mono mt-0.5">{node.mid}</p>
              </button>
            ))}
          </div>
        )}

        {!loading && nodes.length === 0 && (
          <p className="text-sm text-slate-400 px-4 py-8 text-center">No devices registered</p>
        )}
      </div>
    </div>
  );
}
