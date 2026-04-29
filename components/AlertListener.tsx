'use client';

import { useEffect, useRef } from 'react';
import { useAlerts } from '@/internal/domain/service/http/alert';
import { formatDateTime } from '@/internal/utils/format';

function playAlarmTone(ctx: AudioContext) {
  const frequencies = [880, 1100, 880, 1100];
  let t = ctx.currentTime;
  frequencies.forEach((freq) => {
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.frequency.value = freq;
    osc.type = 'square';
    gain.gain.setValueAtTime(0.3, t);
    gain.gain.exponentialRampToValueAtTime(0.001, t + 0.18);
    osc.start(t);
    osc.stop(t + 0.18);
    t += 0.22;
  });
}

export default function AlertListener() {
  const { alerts, acknowledge } = useAlerts(5000);
  const audioCtxRef = useRef<AudioContext | null>(null);
  const alarmIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const prevCountRef = useRef(0);

  // Play alarm whenever there are unacknowledged alerts
  useEffect(() => {
    if (alerts.length === 0) {
      if (alarmIntervalRef.current) {
        clearInterval(alarmIntervalRef.current);
        alarmIntervalRef.current = null;
      }
      prevCountRef.current = 0;
      return;
    }

    // New alert arrived — ensure AudioContext exists (requires user gesture first touch)
    if (!audioCtxRef.current) {
      try {
        audioCtxRef.current = new AudioContext();
      } catch {
        return;
      }
    }

    // Play immediately on first alert or when count increases
    if (alerts.length > prevCountRef.current) {
      playAlarmTone(audioCtxRef.current);
    }
    prevCountRef.current = alerts.length;

    // Repeat alarm every 8 seconds while unacknowledged alerts exist
    if (!alarmIntervalRef.current) {
      alarmIntervalRef.current = setInterval(() => {
        if (audioCtxRef.current) playAlarmTone(audioCtxRef.current);
      }, 8000);
    }

    return () => {};
  }, [alerts]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (alarmIntervalRef.current) clearInterval(alarmIntervalRef.current);
    };
  }, []);

  if (alerts.length === 0) return null;

  return (
    <div className="fixed top-0 left-0 right-0 z-[2000] flex flex-col gap-2 p-3 pointer-events-none">
      {alerts.map((alert) => (
        <div
          key={alert.id}
          className="pointer-events-auto flex items-start gap-3 bg-red-600 text-white px-4 py-3 rounded-xl shadow-2xl animate-pulse"
        >
          <span className="text-2xl shrink-0 mt-0.5">🚨</span>
          <div className="flex-1 min-w-0">
            <p className="font-semibold text-sm">Emergency Alert #{alert.id}</p>
            <p className="text-xs text-red-100 truncate">{alert.message || 'Emergency! Help needed!'}</p>
            <p className="text-[10px] text-red-200 mt-0.5">{formatDateTime(alert.created_at)}</p>
          </div>
          <button
            onClick={() => acknowledge(alert.id)}
            className="shrink-0 bg-white/20 hover:bg-white/30 rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors"
          >
            Acknowledge
          </button>
        </div>
      ))}
    </div>
  );
}
