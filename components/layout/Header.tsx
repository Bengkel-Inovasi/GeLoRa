'use client';

import { useMe } from '../../internal/domain/service/http/user';
import { useAuth } from '../../internal/domain/service/auth';
import Badge from '../ui/Badge';

interface Props {
  title: string;
}

export default function Header({ title }: Props) {
  const { user } = useMe();
  const { signOut } = useAuth();

  return (
    <header className="h-14 bg-white border-b border-slate-200 flex items-center justify-between px-4 md:px-6 shrink-0">
      <h1 className="text-base font-semibold text-slate-800">{title}</h1>

      <div className="flex items-center gap-2">
        {user && (
          <>
            <span className="hidden sm:inline text-sm text-slate-600 truncate max-w-[120px]">{user.name}</span>
            <Badge label={user.role} variant={user.role === 'super' ? 'green' : user.role === 'admin' ? 'blue' : 'gray'} />
          </>
        )}
        {/* Logout button — only shown on mobile (sidebar handles desktop) */}
        <button
          onClick={signOut}
          className="md:hidden flex items-center justify-center w-9 h-9 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors"
          aria-label="Sign out"
          title="Sign out"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-5 h-5">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <polyline points="16 17 21 12 16 7" />
            <line x1="21" y1="12" x2="9" y2="12" />
          </svg>
        </button>
      </div>
    </header>
  );
}
