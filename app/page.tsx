'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { tokenStore } from '../internal/utils/token';

export default function RootPage() {
  const router = useRouter();
  useEffect(() => {
    if (tokenStore.hasToken()) {
      router.replace('/dashboard');
    } else {
      router.replace('/sign-in');
    }
  }, [router]);
  return null;
}
