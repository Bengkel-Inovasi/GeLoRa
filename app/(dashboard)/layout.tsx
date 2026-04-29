'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { tokenStore } from '../../internal/utils/token';
import Sidebar from '../../components/layout/Sidebar';
import SOSButton from '../../components/ui/SOSButton';
import AlertListener from '../../components/AlertListener';

const allNavItems = [
  { href: '/dashboard', label: 'Map', icon: '🗺️', operatorOnly: false },
  { href: '/nodes', label: 'Climbers', icon: '🧗', operatorOnly: true },
  { href: '/sessions', label: 'History', icon: '⛰️', operatorOnly: true },
  { href: '/users', label: 'Users', icon: '👥', operatorOnly: true },
];

// Routes that viewers (client role) cannot access
const OPERATOR_ROUTES = ['/nodes', '/sessions', '/users'];

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [authed, setAuthed] = useState(false);
  const [isOperator, setIsOperator] = useState(false);

  useEffect(() => {
    if (!tokenStore.hasToken()) {
      router.replace('/sign-in');
      return;
    }

    const role = tokenStore.getRole();
    const operator = role === 'admin' || role === 'super';
    setIsOperator(operator);
    setAuthed(true);

    // Redirect viewers away from operator-only pages
    if (!operator && OPERATOR_ROUTES.some((r) => pathname.startsWith(r))) {
      router.replace('/dashboard');
    }
  }, [router, pathname]);

  if (!authed) return null;

  const navItems = allNavItems.filter((item) => !item.operatorOnly || isOperator);

  return (
    <div className="flex h-screen overflow-hidden bg-slate-50">
      {isOperator && <AlertListener />}
      <Sidebar isOperator={isOperator} />
      <div className="flex flex-col flex-1 overflow-hidden pb-16 md:pb-0">
        {children}
      </div>
      <SOSButton />

      {/* Mobile bottom navigation */}
      <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-slate-200 flex md:hidden z-50"
        style={{ height: '64px' }}>
        {navItems.map((item) => {
          const isActive =
            pathname === item.href ||
            (item.href !== '/dashboard' && pathname.startsWith(item.href));
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex-1 flex flex-col items-center justify-center gap-0.5 transition-colors ${
                isActive ? 'text-emerald-600' : 'text-slate-400'
              }`}
            >
              <span className="text-xl leading-none">{item.icon}</span>
              <span className="text-[10px] font-medium">{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
