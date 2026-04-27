'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '../../internal/domain/service/auth';

const allNavItems = [
  { href: '/dashboard', label: 'Live Tracker', icon: '🗺️', operatorOnly: false },
  { href: '/nodes', label: 'Climbers', icon: '🧗', operatorOnly: true },
  { href: '/sessions', label: 'Climb History', icon: '⛰️', operatorOnly: true },
  { href: '/users', label: 'Users', icon: '👥', operatorOnly: true },
];

interface Props {
  isOperator: boolean;
}

export default function Sidebar({ isOperator }: Props) {
  const pathname = usePathname();
  const { signOut } = useAuth();

  const navItems = allNavItems.filter((item) => !item.operatorOnly || isOperator);

  return (
    <aside className="hidden md:flex flex-col w-56 bg-slate-900 text-white h-full shrink-0">
      <div className="px-5 py-5 border-b border-slate-700">
        <p className="text-lg font-bold tracking-tight">GeLoRa</p>
        <p className="text-xs text-slate-400 mt-0.5">Mountain Tracker</p>
      </div>

      <nav className="flex-1 py-4 flex flex-col gap-1 px-3">
        {navItems.map((item) => {
          const isActive = pathname === item.href || (item.href !== '/dashboard' && pathname.startsWith(item.href));
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                isActive ? 'bg-emerald-600 text-white' : 'text-slate-300 hover:bg-slate-800'
              }`}
            >
              <span>{item.icon}</span>
              {item.label}
            </Link>
          );
        })}
      </nav>

      <div className="px-3 py-4 border-t border-slate-700">
        <button
          onClick={signOut}
          className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition-colors"
        >
          <span>🚪</span>
          Sign Out
        </button>
      </div>
    </aside>
  );
}
