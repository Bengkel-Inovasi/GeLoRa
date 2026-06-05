'use client';

import { useState } from 'react';
import { useNodes } from '@/internal/domain/service/http/node';
import { useSessions } from '@/internal/domain/service/http/session';
import { useMe } from '@/internal/domain/service/http/user';
import { nodeRepository } from '@/internal/adapters/outbound/repository/node';
import { sessionRepository } from '@/internal/adapters/outbound/repository/session';
import Header from '@/components/layout/Header';
import Card from '@/components/ui/Card';
import Badge from '@/components/ui/Badge';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Spinner from '@/components/ui/Spinner';
import { formatTime } from '@/internal/utils/format';
import type { Session } from '@/internal/domain/model/session';

export default function ClimbersPage() {
  const { nodes, loading, error, refetch: refetchNodes } = useNodes();
  const { sessions, refetch: refetchSessions } = useSessions();
  const { user } = useMe();

  const [showModal, setShowModal] = useState(false);
  const [mid, setMid] = useState('');
  const [climberName, setClimberName] = useState('');
  const [registering, setRegistering] = useState(false);
  const [registerError, setRegisterError] = useState<string | null>(null);
  const [actionId, setActionId] = useState<number | null>(null);

  const isAdmin = user?.role === 'admin' || user?.role === 'super';

  const activeSessionByNode = new Map<number, Session>();
  for (const s of sessions) {
    if (s.ended_at == null && s.node_id != null) {
      activeSessionByNode.set(s.node_id, s);
    }
  }

  function openModal() {
    setMid('');
    setClimberName('');
    setRegisterError(null);
    setShowModal(true);
  }

  async function handleRegister(e: React.FormEvent) {
    e.preventDefault();
    if (!mid.trim() || !climberName.trim()) return;
    setRegistering(true);
    setRegisterError(null);
    try {
      await nodeRepository.register(mid.trim(), climberName.trim());
      setShowModal(false);
      await Promise.all([refetchNodes(), refetchSessions()]);
    } catch (err) {
      setRegisterError(err instanceof Error ? err.message : 'Registration failed. Check the LoRa address and try again.');
    } finally {
      setRegistering(false);
    }
  }

  async function handleEndClimb(nodeId: number) {
    setActionId(nodeId);
    try {
      await sessionRepository.end({ node_id: nodeId });
      await refetchSessions();
    } finally {
      setActionId(null);
    }
  }

  async function handleRemove(id: number, name: string, mid: string) {
    if (!confirm(`Remove ${name} (${mid}) from the system?\n\nThe wristband will be free to assign to a new climber.`)) return;
    setActionId(id);
    try {
      await nodeRepository.delete(id);
      await Promise.all([refetchNodes(), refetchSessions()]);
    } finally {
      setActionId(null);
    }
  }

  return (
    <>
      <Header title="Climbers" />

      <div className="flex-1 overflow-y-auto p-4 md:p-6">
        {isAdmin && (
          <div className="flex items-center justify-between mb-5">
            <p className="text-sm text-slate-500">
              {nodes.length} wristband{nodes.length !== 1 ? 's' : ''} registered
            </p>
            <Button onClick={openModal} size="sm">+ Register Climber</Button>
          </div>
        )}

        {loading && (
          <div className="flex justify-center py-20">
            <Spinner size={32} />
          </div>
        )}
        {error && <p className="text-red-500 text-sm mb-4">{error}</p>}

        {!loading && nodes.length === 0 && (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <p className="text-5xl mb-4">📡</p>
            <p className="font-semibold text-slate-700 text-lg mb-2">No climbers registered yet</p>
            <p className="text-sm text-slate-400 mb-8 max-w-sm">
              Register a climber to link their wristband and start tracking their location and health data in real time.
            </p>
            {isAdmin && <Button onClick={openModal}>+ Register Climber</Button>}
          </div>
        )}

        <div className="grid gap-6">
          {/* ── New Devices Section ── */}
          {nodes.some(n => !n.is_validated) && (
            <section>
              <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3 px-1">
                New Devices Detected (Unvalidated)
              </h3>
              <div className="grid gap-3">
                {nodes.filter(n => !n.is_validated).map((node) => {
                  const busy = actionId === node.id;
                  return (
                    <Card key={node.id} className="border-l-4 border-l-amber-400">
                      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-0.5">
                            <p className="font-mono font-bold text-slate-800">{node.mid}</p>
                            <Badge label="New Device" variant="red" />
                          </div>
                          <p className="text-xs text-slate-500 italic">
                            Waiting for admin to assign a climber name and validate.
                          </p>
                        </div>

                        {isAdmin && (
                          <div className="flex gap-2 shrink-0">
                            <Button size="sm" onClick={() => {
                              setMid(node.mid);
                              setClimberName('');
                              setShowModal(true);
                            }}>
                              Validate &amp; Alias
                            </Button>
                            <Button size="sm" variant="danger" loading={busy} onClick={() => handleRemove(node.id, node.name, node.mid)}>
                              Discard
                            </Button>
                          </div>
                        )}
                      </div>
                    </Card>
                  );
                })}
              </div>
            </section>
          )}

          {/* ── Registered Climbers Section ── */}
          <section>
            <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3 px-1">
              Registered Climbers
            </h3>
            <div className="grid gap-3">
              {nodes.filter(n => n.is_validated).map((node) => {
                const session = activeSessionByNode.get(node.id);
                const isTracking = !!session;
                const busy = actionId === node.id;
                return (
                  <Card key={node.id}>
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 flex-wrap mb-0.5">
                          <p className="font-semibold text-slate-800">{node.name}</p>
                          <Badge
                            label={isTracking ? '● Tracking Active' : '○ No Active Session'}
                            variant={isTracking ? 'green' : 'gray'}
                          />
                        </div>
                        <p className="text-xs text-slate-400 font-mono">Wristband: {node.mid}</p>
                        {session && (
                          <p className="text-xs text-slate-400 mt-0.5">
                            Session #{session.id} · Started {formatTime(session.started_at)}
                          </p>
                        )}
                      </div>

                      {isAdmin && (
                        <div className="flex gap-2 shrink-0">
                          {isTracking && (
                            <Button size="sm" variant="secondary" loading={busy} onClick={() => handleEndClimb(node.id)}>
                              End Climb
                            </Button>
                          )}
                          <Button size="sm" variant="danger" loading={busy} onClick={() => handleRemove(node.id, node.name, node.mid)}>
                            Remove
                          </Button>
                        </div>
                      )}
                    </div>
                  </Card>
                );
              })}
              {nodes.filter(n => n.is_validated).length === 0 && (
                <p className="text-sm text-slate-400 italic py-4 px-1">No climbers validated yet.</p>
              )}
            </div>
          </section>
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-4">
          <div className="bg-white w-full sm:max-w-md sm:rounded-2xl rounded-t-2xl shadow-2xl">
            {/* Drag handle for mobile */}
            <div className="flex justify-center pt-3 pb-1 sm:hidden">
              <div className="w-10 h-1 bg-slate-300 rounded-full" />
            </div>

            <div className="flex items-start justify-between px-6 pt-4 pb-4 border-b border-slate-100">
              <div>
                <h2 className="text-lg font-semibold text-slate-800">Register New Climber</h2>
                <p className="text-xs text-slate-400 mt-0.5">Links a wristband to a climber and starts tracking immediately</p>
              </div>
              <button
                onClick={() => setShowModal(false)}
                className="text-slate-400 hover:text-slate-600 text-2xl leading-none ml-4 -mt-1"
              >
                &times;
              </button>
            </div>

            <form onSubmit={handleRegister} className="p-6 flex flex-col gap-5">
              <div className="flex flex-col gap-1">
                <Input
                  label="LoRa Address (Wristband ID)"
                  placeholder="e.g. PENDAKI_01"
                  value={mid}
                  onChange={(e) => setMid(e.target.value)}
                  required
                />
                <p className="text-xs text-slate-400 pl-0.5">The hardware ID shown on the back of the wristband</p>
              </div>

              <Input
                label="Climber Name"
                placeholder="e.g. Ahmad Fauzi"
                value={climberName}
                onChange={(e) => setClimberName(e.target.value)}
                required
              />

              {registerError && (
                <p className="text-sm text-red-500 bg-red-50 px-3 py-2 rounded-lg">{registerError}</p>
              )}

              <div className="flex gap-3 pt-1 pb-safe">
                <Button type="button" variant="secondary" className="flex-1" onClick={() => setShowModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" loading={registering} className="flex-1">
                  Register &amp; Start Tracking
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
