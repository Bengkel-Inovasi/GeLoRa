import { http, HttpResponse } from 'msw';
import { API_BASE } from '../config/constants';
import { mockNodes, mockSessions, mockUsers, getMockRecords } from './data';

const pagination = (data: unknown[]) => ({
  page: 1,
  limit: 100,
  total: data.length,
});

interface MockUser {
  id: number;
  name: string;
  username: string;
  password: string;
  role: string;
  bio: string;
}

const MOCK_ACCOUNTS: MockUser[] = [
  { id: 1, name: 'Super Admin',    username: 'superadmin', password: 'super@gelora',    role: 'super',  bio: '' },
  { id: 2, name: 'Default Admin',  username: 'admin',      password: 'admin@gelora',    role: 'admin',  bio: '' },
  { id: 3, name: 'Default Client', username: 'client',     password: 'client@gelora',   role: 'client', bio: '' },
  { id: 5, name: 'Operator',       username: 'Operator',   password: 'SuperOperator67', role: 'admin',  bio: 'Basecamp operator' },
];

let currentMockUser: MockUser = MOCK_ACCOUNTS[0];

export const handlers = [
  http.post(`${API_BASE}/auth/sign-in`, async ({ request }) => {
    const body = await request.json() as { username?: string; password?: string };
    const account = MOCK_ACCOUNTS.find(
      (a) => a.username === body.username && a.password === body.password,
    );
    if (!account) {
      return HttpResponse.json({ error: 'Invalid username or password' }, { status: 401 });
    }
    currentMockUser = account;
    return HttpResponse.json({
      access_token: `mock-token-${account.username}`,
      refresh_token: `mock-refresh-${account.username}`,
    });
  }),

  http.post(`${API_BASE}/auth/refresh`, () =>
    HttpResponse.json({
      access_token: `mock-token-${currentMockUser.username}`,
      refresh_token: `mock-refresh-${currentMockUser.username}`,
    }),
  ),

  http.get(`${API_BASE}/users/me`, () =>
    HttpResponse.json({
      id: currentMockUser.id,
      name: currentMockUser.name,
      username: currentMockUser.username,
      role: currentMockUser.role,
      bio: currentMockUser.bio,
      created_at: new Date().toISOString(),
      updated_at: null,
    }),
  ),

  http.get(`${API_BASE}/users`, () =>
    HttpResponse.json({ data: mockUsers, pagination: pagination(mockUsers) }),
  ),

  http.get(`${API_BASE}/nodes`, () =>
    HttpResponse.json({ data: mockNodes, pagination: pagination(mockNodes) }),
  ),

  http.get(`${API_BASE}/nodes/:id`, ({ params }) => {
    const node = mockNodes.find((n) => n.id === Number(params.id));
    if (!node) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json({ ...node, description: 'Mock device', validated_by: null, created_at: new Date().toISOString(), updated_at: null });
  }),

  http.post(`${API_BASE}/nodes/register`, () =>
    HttpResponse.json({ node_id: Math.floor(Math.random() * 900 + 100), session_id: Math.floor(Math.random() * 900 + 100) }, { status: 201 }),
  ),
  http.patch(`${API_BASE}/nodes/:id`, () => new HttpResponse(null, { status: 204 })),
  http.put(`${API_BASE}/nodes/:id/validate`, () => new HttpResponse(null, { status: 204 })),
  http.delete(`${API_BASE}/nodes/:id`, () => new HttpResponse(null, { status: 204 })),

  http.get(`${API_BASE}/sessions`, () =>
    HttpResponse.json({ data: mockSessions, pagination: pagination(mockSessions) }),
  ),

  http.get(`${API_BASE}/sessions/:id`, ({ params }) => {
    const session = mockSessions.find((s) => s.id === Number(params.id));
    if (!session) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(session);
  }),

  http.post(`${API_BASE}/sessions`, () => HttpResponse.json({ id: 99 })),
  http.put(`${API_BASE}/sessions/end`, () => new HttpResponse(null, { status: 204 })),
  http.delete(`${API_BASE}/sessions/:id`, () => new HttpResponse(null, { status: 204 })),

  http.get(`${API_BASE}/records`, ({ request }) => {
    const url = new URL(request.url);
    const nodeId = url.searchParams.get('node_id');
    const records = getMockRecords();
    const filtered = nodeId
      ? records.filter((r) => r.session_id === Number(nodeId))
      : records;
    return HttpResponse.json({ data: filtered });
  }),
];
