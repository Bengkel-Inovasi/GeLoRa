import type { Node, NodeListItem } from '@/internal/domain/model/node';
import type { PaginationMeta } from '@/internal/domain/model/error';

export interface INodeRepository {
  list(params?: { page?: number; limit?: number }): Promise<{ data: NodeListItem[]; pagination: PaginationMeta }>;
  getById(id: number): Promise<Node>;
  patch(id: number, payload: { name?: string; description?: string }): Promise<void>;
  validate(id: number, validatedBy: number): Promise<void>;
  delete(id: number): Promise<void>;
  register(mid: string, name: string): Promise<{ node_id: number; session_id: number }>;
}
