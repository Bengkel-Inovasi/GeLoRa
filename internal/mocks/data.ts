import type { NodeListItem } from '../domain/model/node';
import type { Session } from '../domain/model/session';
import type { Record } from '../domain/model/record';
import type { UserListItem } from '../domain/model/user';

export const mockNodes: NodeListItem[] = [
  { id: 1, mid: 'gelora-a1b2c3', name: 'Budi Santoso', is_validated: true },
  { id: 2, mid: 'gelora-d4e5f6', name: 'Siti Rahayu', is_validated: true },
  { id: 3, mid: 'gelora-g7h8i9', name: 'Ahmad Fauzi', is_validated: true },
  { id: 4, mid: 'gelora-j1k2l3', name: 'Dewi Kusuma', is_validated: false },
];

export const mockSessions: Session[] = [
  { id: 1, user_id: 2, node_id: 1, started_at: new Date(Date.now() - 3600000).toISOString(), ended_at: null },
  { id: 2, user_id: 3, node_id: 2, started_at: new Date(Date.now() - 7200000).toISOString(), ended_at: null },
  { id: 3, user_id: 4, node_id: 3, started_at: new Date(Date.now() - 1800000).toISOString(), ended_at: null },
  { id: 4, user_id: 2, node_id: 1, started_at: new Date(Date.now() - 86400000).toISOString(), ended_at: new Date(Date.now() - 72000000).toISOString() },
];

export const mockUsers: UserListItem[] = [
  { id: 1, name: 'Super Admin', username: 'superadmin', role: 'super' },
  { id: 2, name: 'Eko Prasetyo', username: 'eko.prasetyo', role: 'admin' },
  { id: 3, name: 'Wati Indah', username: 'wati.indah', role: 'client' },
  { id: 4, name: 'Roni Hidayat', username: 'roni.hidayat', role: 'client' },
  { id: 5, name: 'Operator', username: 'Operator', role: 'admin' },
];

const BASE_COORDS: [number, number][] = [
  [-7.9167, 110.3333],  // Merapi, Java
  [-8.3405, 115.508],   // Agung, Bali
  [-8.9408, 116.4660],  // Rinjani, Lombok
];

let tick = 0;

export function getMockRecords(): Record[] {
  tick++;
  const now = new Date().toISOString();
  return [1, 2, 3].map((nodeId, i) => {
    const [baseLat, baseLon] = BASE_COORDS[i];
    const jitter = () => (Math.random() - 0.5) * 0.002 * tick * 0.1;
    return {
      id: tick * 10 + nodeId,
      session_id: nodeId,
      time: now,
      heart_rate: 72 + Math.round(Math.random() * 30 - 5),
      body_temperature: 36.2 + Math.random() * 1.5,
      latitude: baseLat + jitter(),
      longitude: baseLon + jitter(),
    };
  });
}
