'use client';

import { useMe } from '../../internal/domain/service/http/user';
import Badge from '../ui/Badge';

interface Props {
  title: string;
}

export default function Header({ title }: Props) {
  const { user } = useMe();

  return (
    <header className="h-14 bg-white border-b border-slate-200 flex items-center justify-between px-4 md:px-6 shrink-0">
      <h1 className="text-base font-semibold text-slate-800">{title}</h1>
      {user && (
        <div className="flex items-center gap-2 text-sm text-slate-600">
          <span className="hidden sm:inline truncate max-w-[120px]">{user.name}</span>
          <Badge label={user.role} variant={user.role === 'super' ? 'green' : user.role === 'admin' ? 'blue' : 'gray'} />
        </div>
      )}
    </header>
  );
}
