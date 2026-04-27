'use client';

import { useUsers, useMe } from '../../../internal/domain/service/http/user';
import { userRepository } from '../../../internal/adapters/outbound/repository/user';
import { authRepository } from '../../../internal/adapters/outbound/repository/auth';
import Header from '../../../components/layout/Header';
import Card from '../../../components/ui/Card';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Input from '../../../components/ui/Input';
import Spinner from '../../../components/ui/Spinner';
import { useState } from 'react';
import type { UserRole } from '../../../internal/domain/model/user';

const ROLES: UserRole[] = ['client', 'admin', 'super'];

type NewAccountRole = 'viewer' | 'operator';

function decodeJwtId(token: string): number | null {
  try {
    const b64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    const payload = JSON.parse(atob(b64));
    return typeof payload.Id === 'number' ? payload.Id : null;
  } catch {
    return null;
  }
}

export default function UsersPage() {
  const { users, loading, error, refetch } = useUsers();
  const { user: me } = useMe();
  const [actionId, setActionId] = useState<number | null>(null);
  const isSuper = me?.role === 'super';

  // Create account state
  const [showCreate, setShowCreate] = useState(false);
  const [newName, setNewName] = useState('');
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newRole, setNewRole] = useState<NewAccountRole>('viewer');
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);

  function openCreate() {
    setNewName(''); setNewUsername(''); setNewPassword('');
    setNewRole('viewer'); setCreateError(null);
    setShowCreate(true);
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!newName.trim() || !newUsername.trim() || !newPassword.trim()) return;
    setCreating(true);
    setCreateError(null);
    try {
      // Create account via sign-up (returns tokens with the new user's JWT)
      const tokens = await authRepository.signUp({
        name: newName.trim(),
        username: newUsername.trim(),
        password: newPassword.trim(),
      });
      // If operator role selected, decode the new user's ID from the JWT and promote
      if (newRole === 'operator') {
        const newUserId = decodeJwtId(tokens.access_token);
        if (newUserId) {
          await userRepository.updateRoleById(newUserId, 'admin');
        }
      }
      setShowCreate(false);
      await refetch();
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : 'Failed to create account. Username may already be taken.');
    } finally {
      setCreating(false);
    }
  }

  async function handleRoleChange(id: number, role: UserRole) {
    setActionId(id);
    try {
      await userRepository.updateRoleById(id, role);
      await refetch();
    } finally {
      setActionId(null);
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Delete this user?')) return;
    setActionId(id);
    try {
      await userRepository.delete(id);
      await refetch();
    } finally {
      setActionId(null);
    }
  }

  return (
    <>
      <Header title="Users" />
      <div className="flex-1 overflow-y-auto p-4 md:p-6">

        {isSuper && (
          <div className="flex items-center justify-between mb-5">
            <p className="text-sm text-slate-500">{users.length} account{users.length !== 1 ? 's' : ''}</p>
            <Button size="sm" onClick={openCreate}>+ Create Account</Button>
          </div>
        )}

        {loading && <div className="flex justify-center py-16"><Spinner size={32} /></div>}
        {error && <p className="text-red-500 text-sm">{error}</p>}

        <div className="grid gap-3">
          {users.map((u) => (
            <Card key={u.id}>
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div className="min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <p className="font-semibold text-slate-800">{u.name}</p>
                    <Badge
                      label={u.role === 'client' ? 'Viewer' : u.role === 'admin' ? 'Operator' : 'Super'}
                      variant={u.role === 'super' ? 'green' : u.role === 'admin' ? 'blue' : 'gray'}
                    />
                    {me?.id === u.id && <Badge label="You" variant="gray" />}
                  </div>
                  <p className="text-xs text-slate-400 font-mono mt-0.5">@{u.username}</p>
                </div>
                {isSuper && me?.id !== u.id && (
                  <div className="flex items-center gap-2 shrink-0">
                    <select
                      className="rounded-lg border border-slate-300 px-2 py-1.5 text-sm flex-1 sm:flex-none"
                      value={u.role}
                      onChange={(e) => handleRoleChange(u.id, e.target.value as UserRole)}
                      disabled={actionId === u.id}
                    >
                      {ROLES.map((r) => (
                        <option key={r} value={r}>
                          {r === 'client' ? 'Viewer' : r === 'admin' ? 'Operator' : 'Super'}
                        </option>
                      ))}
                    </select>
                    <Button size="sm" variant="danger" loading={actionId === u.id} onClick={() => handleDelete(u.id)}>
                      Delete
                    </Button>
                  </div>
                )}
              </div>
            </Card>
          ))}
          {!loading && users.length === 0 && (
            <p className="text-slate-400 text-sm text-center py-16">No users found</p>
          )}
        </div>
      </div>

      {/* Create Account modal */}
      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50">
          <div className="bg-white w-full sm:max-w-md sm:rounded-2xl rounded-t-2xl shadow-2xl">
            <div className="flex justify-center pt-3 pb-1 sm:hidden">
              <div className="w-10 h-1 bg-slate-300 rounded-full" />
            </div>

            <div className="flex items-start justify-between px-6 pt-4 pb-4 border-b border-slate-100">
              <div>
                <h2 className="text-lg font-semibold text-slate-800">Create Account</h2>
                <p className="text-xs text-slate-400 mt-0.5">Viewer: map only · Operator: full control</p>
              </div>
              <button onClick={() => setShowCreate(false)} className="text-slate-400 hover:text-slate-600 text-2xl leading-none ml-4 -mt-1">&times;</button>
            </div>

            <form onSubmit={handleCreate} className="p-6 flex flex-col gap-4">
              <Input label="Full Name" placeholder="e.g. Siti Rahayu" value={newName} onChange={(e) => setNewName(e.target.value)} required />
              <Input label="Username (min. 3 characters)" placeholder="e.g. siti" value={newUsername} onChange={(e) => setNewUsername(e.target.value)} required />
              <Input label="Password" type="password" placeholder="Min. 8 characters" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required />

              {/* Role picker */}
              <div className="flex flex-col gap-1.5">
                <p className="text-sm font-medium text-slate-700">Account Type</p>
                <div className="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    onClick={() => setNewRole('viewer')}
                    className={`rounded-xl border-2 px-4 py-3 text-left transition-colors ${
                      newRole === 'viewer' ? 'border-emerald-500 bg-emerald-50' : 'border-slate-200 hover:border-slate-300'
                    }`}
                  >
                    <p className="font-semibold text-sm text-slate-800">👁️ Viewer</p>
                    <p className="text-xs text-slate-500 mt-0.5">Map & location only</p>
                  </button>
                  <button
                    type="button"
                    onClick={() => setNewRole('operator')}
                    className={`rounded-xl border-2 px-4 py-3 text-left transition-colors ${
                      newRole === 'operator' ? 'border-blue-500 bg-blue-50' : 'border-slate-200 hover:border-slate-300'
                    }`}
                  >
                    <p className="font-semibold text-sm text-slate-800">🛠️ Operator</p>
                    <p className="text-xs text-slate-500 mt-0.5">Full management access</p>
                  </button>
                </div>
              </div>

              {createError && (
                <p className="text-sm text-red-500 bg-red-50 px-3 py-2 rounded-lg">{createError}</p>
              )}

              <div className="flex gap-3 pt-1">
                <Button type="button" variant="secondary" className="flex-1" onClick={() => setShowCreate(false)}>Cancel</Button>
                <Button type="submit" loading={creating} className="flex-1">Create Account</Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
