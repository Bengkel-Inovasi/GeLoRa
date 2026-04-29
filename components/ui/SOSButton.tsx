'use client';

import { useState } from 'react';
import { sendAlert } from '@/internal/domain/service/http/alert';
import { useNodes } from '@/internal/domain/service/http/node';

export default function SOSButton() {
  const [open, setOpen] = useState(false);
  const [selectedNodeId, setSelectedNodeId] = useState<number | null>(null);
  const [message, setMessage] = useState('');
  const [sending, setSending] = useState(false);
  const [sent, setSent] = useState(false);
  const { nodes } = useNodes();

  function handleOpen() {
    setOpen(true);
    setSent(false);
    setSelectedNodeId(null);
    setMessage('');
  }

  async function handleSend() {
    if (sending) return;
    setSending(true);
    try {
      await sendAlert(selectedNodeId, message.trim() || 'Emergency! I need help!');
      setSent(true);
      setTimeout(() => {
        setSent(false);
        setOpen(false);
      }, 3000);
    } catch {
      alert('Failed to send SOS. Check your connection.');
    } finally {
      setSending(false);
    }
  }

  return (
    <>
      <button
        onClick={handleOpen}
        className="fixed bottom-20 right-4 md:bottom-6 md:right-6 z-[1100] w-14 h-14 rounded-full bg-red-600 hover:bg-red-700 active:scale-95 transition-all shadow-lg flex items-center justify-center text-white font-bold text-sm select-none"
        aria-label="Send SOS emergency alert"
      >
        SOS
      </button>

      {open && (
        <div className="fixed inset-0 z-[1200] flex items-end sm:items-center justify-center">
          <div className="absolute inset-0 bg-black/50" onClick={() => setOpen(false)} />
          <div className="relative bg-white rounded-t-2xl sm:rounded-2xl p-6 w-full max-w-sm mx-0 sm:mx-4 shadow-2xl">
            {sent ? (
              <div className="flex flex-col items-center gap-3 py-4">
                <div className="w-14 h-14 rounded-full bg-green-100 flex items-center justify-center text-3xl">✓</div>
                <p className="font-semibold text-slate-800">SOS Sent!</p>
                <p className="text-sm text-slate-500 text-center">Operators have been alerted and will respond.</p>
              </div>
            ) : (
              <>
                <div className="flex items-center gap-3 mb-5">
                  <div className="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center text-xl shrink-0">🚨</div>
                  <div>
                    <p className="font-semibold text-slate-800">Emergency Alert</p>
                    <p className="text-xs text-slate-500">Operators will be notified immediately</p>
                  </div>
                </div>

                {/* Climber selector */}
                <div className="mb-4">
                  <label className="block text-xs font-medium text-slate-600 mb-1.5">
                    Who is in danger? <span className="text-red-500">*</span>
                  </label>
                  <select
                    value={selectedNodeId ?? ''}
                    onChange={(e) => setSelectedNodeId(e.target.value ? Number(e.target.value) : null)}
                    className="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-red-300 text-slate-800"
                  >
                    <option value="">— Select a climber —</option>
                    {nodes.map((node) => (
                      <option key={node.id} value={node.id}>
                        {node.name} · {node.mid}
                      </option>
                    ))}
                  </select>
                  {nodes.length === 0 && (
                    <p className="text-xs text-slate-400 mt-1">No climbers registered yet.</p>
                  )}
                </div>

                {/* Optional message */}
                <div className="mb-5">
                  <label className="block text-xs font-medium text-slate-600 mb-1.5">
                    Additional details (optional)
                  </label>
                  <textarea
                    className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-red-300"
                    rows={2}
                    placeholder="Describe the situation..."
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    maxLength={500}
                  />
                </div>

                <div className="flex gap-2">
                  <button
                    onClick={() => setOpen(false)}
                    className="flex-1 py-2.5 rounded-lg border border-slate-200 text-sm font-medium text-slate-600 hover:bg-slate-50 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleSend}
                    disabled={sending || selectedNodeId === null}
                    className="flex-1 py-2.5 rounded-lg bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-semibold text-white transition-colors"
                  >
                    {sending ? 'Sending...' : 'Send SOS'}
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </>
  );
}
