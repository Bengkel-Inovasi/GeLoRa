'use client';

import { useState } from 'react';
import { useSessions } from '@/internal/domain/service/http/session';
import { useNodes } from '@/internal/domain/service/http/node';
import { sessionRepository } from '@/internal/adapters/outbound/repository/session';
import Header from '@/components/layout/Header';
import Card from '@/components/ui/Card';
import Badge from '@/components/ui/Badge';
import Button from '@/components/ui/Button';
import Spinner from '@/components/ui/Spinner';
import { formatDateTime } from '@/internal/utils/format';

export default function SessionsPage() {
  const { sessions, loading, error, refetch } = useSessions();
  const { nodes } = useNodes();
  const [actionId, setActionId] = useState<number | null>(null);
  const [filter, setFilter] = useState<'all' | 'active' | 'ended'>('all');

  const nodeMap = Object.fromEntries(nodes.map((n) => [n.id, n.name]));

  const filtered = sessions.filter((s) => {
    if (filter === 'active') return s.ended_at == null;
    if (filter === 'ended') return s.ended_at != null;
    return true;
  });

  const activeCount = sessions.filter((s) => s.ended_at == null).length;

  async function handleEnd(id: number) {
    if (!confirm('End this climb session?')) return;
    setActionId(id);
    try {
      await sessionRepository.end({ id });
      await refetch();
    } finally {
      setActionId(null);
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Delete this session record? This cannot be undone.')) return;
    setActionId(id);
    try {
      await sessionRepository.delete(id);
      await refetch();
    } finally {
      setActionId(null);
    }
  }

  return (
    <>
      <Header title="Climb History" />
      <div className="flex-1 overflow-y-auto flex flex-col">

        {/* Filter tabs — sticky so they don't scroll away */}
        <div className="sticky top-0 z-10 bg-white border-b border-slate-100 px-4 md:px-6 py-3 flex items-center gap-3 overflow-x-auto shrink-0">
          <div className="flex gap-1 bg-slate-100 rounded-lg p-1 shrink-0">
            {(['all', 'active', 'ended'] as const).map((f) => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors capitalize whitespace-nowrap ${
                  filter === f ? 'bg-white shadow-sm text-slate-800' : 'text-slate-500 hover:text-slate-700'
                }`}
              >
                {f === 'all'
                  ? `All (${sessions.length})`
                  : f === 'active'
                  ? `Active (${activeCount})`
                  : `Ended (${sessions.length - activeCount})`}
              </button>
            ))}
          </div>
          <p className="hidden sm:block text-xs text-slate-400 shrink-0">
            To start tracking, go to <span className="font-medium text-slate-600">Climbers → Register Climber</span>
          </p>
        </div>

        <div className="p-4 md:p-6 flex flex-col gap-4 md:gap-5 flex-1">
          {loading && <div className="flex justify-center py-16"><Spinner size={32} /></div>}
          {error && <p className="text-red-500 text-sm">{error}</p>}

          {!loading && filtered.length === 0 && (
            <div className="flex flex-col items-center justify-center py-20 text-center">
              <p className="text-4xl mb-3">⛰️</p>
              <p className="font-semibold text-slate-700 mb-1">No sessions found</p>
              <p className="text-sm text-slate-400">
                {filter === 'active' ? 'No active climbs right now.' : 'No session history yet.'}
              </p>
            </div>
          )}

          <div className="grid gap-3">
            {filtered.map((session) => {
              const isActive = session.ended_at == null;
              const busy = actionId === session.id;
              const nodeName = session.node_id ? (nodeMap[session.node_id] ?? `Wristband #${session.node_id}`) : '—';

              return (
                <Card key={session.id}>
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap mb-1">
                        <Badge label={isActive ? '● Active' : 'Ended'} variant={isActive ? 'green' : 'gray'} />
                        <span className="text-xs text-slate-400">Session #{session.id}</span>
                      </div>
                      <p className="font-semibold text-slate-800">{nodeName}</p>
                      <p className="text-xs text-slate-400 mt-0.5">
                        Started: {formatDateTime(session.started_at)}
                        {session.ended_at && (
                          <> · Ended: {formatDateTime(session.ended_at)}</>
                        )}
                      </p>
                    </div>

                    <div className="flex gap-2 shrink-0">
                      {isActive && (
                        <Button size="sm" variant="secondary" loading={busy} onClick={() => handleEnd(session.id)}>
                          End Climb
                        </Button>
                      )}
                      <Button size="sm" variant="danger" loading={busy} onClick={() => handleDelete(session.id)}>
                        Delete
                      </Button>
                    </div>
                  </div>
                </Card>
              );
            })}
          </div>
        </div>
      </div>
    </>
  );
}
