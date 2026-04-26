import type { ISessionRepository } from '../../../ports/outbound/repository/session';
import type { Session } from '../../../domain/model/session';
import type { PaginationMeta } from '../../../domain/model/error';
import { httpClient } from '../client/http';

class SessionRepository implements ISessionRepository {
  list(params?: { page?: number; limit?: number }): Promise<{ data: Session[]; pagination: PaginationMeta }> {
    const q = new URLSearchParams();
    if (params?.page) q.set('page', String(params.page));
    if (params?.limit) q.set('limit', String(params.limit));
    const qs = q.toString() ? `?${q}` : '';
    return httpClient.get(`/sessions${qs}`);
  }

  getById(id: number): Promise<Session> {
    return httpClient.get(`/sessions/${id}`);
  }

  create(userId: number, nodeId: number): Promise<{ id: number }> {
    return httpClient.post('/sessions', { user_id: userId, node_id: nodeId });
  }

  end(payload: { id?: number; user_id?: number; node_id?: number }): Promise<void> {
    return httpClient.put('/sessions/end', payload);
  }

  delete(id: number): Promise<void> {
    return httpClient.delete(`/sessions/${id}`);
  }
}

export const sessionRepository = new SessionRepository();
