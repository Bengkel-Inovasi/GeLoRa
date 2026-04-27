'use client';

import { useUsers, useMe } from '../../../internal/domain/service/http/user';
import { userRepository } from '../../../internal/adapters/outbound/repository/user';
import Header from '../../../components/layout/Header';
import Card from '../../../components/ui/Card';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Spinner from '../../../components/ui/Spinner';
import { useState } from 'react';
import type { UserRole } from '../../../internal/domain/model/user';

const ROLES: UserRole[] = ['client', 'admin', 'super'];

export default function UsersPage() {
  const { users, loading, error, refetch } = useUsers();
  const { user: me } = useMe();
  const [actionId, setActionId] = useState<number | null>(null);
  const isSuper = me?.role === 'super';

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
                      label={u.role}
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
                      {ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
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
    </>
  );
}
