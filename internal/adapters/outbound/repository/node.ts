import type { INodeRepository } from '../../../ports/outbound/repository/node';
import type { Node, NodeListItem } from '../../../domain/model/node';
import type { PaginationMeta } from '../../../domain/model/error';
import { httpClient } from '../client/http';

class NodeRepository implements INodeRepository {
  list(params?: { page?: number; limit?: number }): Promise<{ data: NodeListItem[]; pagination: PaginationMeta }> {
    const q = new URLSearchParams();
    if (params?.page) q.set('page', String(params.page));
    if (params?.limit) q.set('limit', String(params.limit));
    const qs = q.toString() ? `?${q}` : '';
    return httpClient.get(`/nodes${qs}`);
  }

  getById(id: number): Promise<Node> {
    return httpClient.get(`/nodes/${id}`);
  }

  patch(id: number, payload: { name?: string; description?: string }): Promise<void> {
    return httpClient.patch(`/nodes/${id}`, payload);
  }

  validate(id: number, validatedBy: number): Promise<void> {
    return httpClient.put(`/nodes/${id}/validate`, { validated_by: validatedBy });
  }

  delete(id: number): Promise<void> {
    return httpClient.delete(`/nodes/${id}`);
  }

  register(mid: string, name: string): Promise<{ node_id: number; session_id: number }> {
    return httpClient.post('/nodes/register', { mid, name });
  }
}

export const nodeRepository = new NodeRepository();
