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
  const [showPassword, setShowPassword] = useState(false);
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
              
              <div className="relative flex flex-col gap-1">
                <Input
                  label="Password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Min. 8 characters"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="pr-10"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-[32px] text-slate-400 hover:text-slate-600 focus:outline-none"
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                >
                  {showPassword ? (
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="h-5 w-5">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                    </svg>
                  ) : (
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="h-5 w-5">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                      <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                  )}
                </button>
              </div>

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
