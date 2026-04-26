import type { Session } from '@/internal/domain/model/session';
import type { PaginationMeta } from '@/internal/domain/model/error';

export interface ISessionRepository {
  list(params?: { page?: number; limit?: number }): Promise<{ data: Session[]; pagination: PaginationMeta }>;
  getById(id: number): Promise<Session>;
  create(userId: number, nodeId: number): Promise<{ id: number }>;
  end(payload: { id?: number; user_id?: number; node_id?: number }): Promise<void>;
  delete(id: number): Promise<void>;
}
