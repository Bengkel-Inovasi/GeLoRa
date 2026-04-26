import type { IUserRepository } from '../../../ports/outbound/repository/user';
import type { User, UserListItem, UserRole } from '../../../domain/model/user';
import type { PaginationMeta } from '../../../domain/model/error';
import { httpClient } from '../client/http';

class UserRepository implements IUserRepository {
  me(): Promise<User> {
    return httpClient.get('/users/me');
  }

  list(params?: { page?: number; limit?: number }): Promise<{ data: UserListItem[]; pagination: PaginationMeta }> {
    const q = new URLSearchParams();
    if (params?.page) q.set('page', String(params.page));
    if (params?.limit) q.set('limit', String(params.limit));
    const qs = q.toString() ? `?${q}` : '';
    return httpClient.get(`/users${qs}`);
  }

  getById(id: number): Promise<User> {
    return httpClient.get(`/users/${id}`);
  }

  patchMe(payload: { name?: string; bio?: string }): Promise<void> {
    return httpClient.patch('/users/me', payload);
  }

  patchById(id: number, payload: { name?: string; bio?: string }): Promise<void> {
    return httpClient.patch(`/users/${id}`, payload);
  }

  updatePasswordMe(oldPassword: string, newPassword: string): Promise<void> {
    return httpClient.put('/users/me/password', { old_password: oldPassword, new_password: newPassword });
  }

  updatePasswordById(id: number, newPassword: string): Promise<void> {
    return httpClient.put(`/users/${id}/password`, { new_password: newPassword });
  }

  updateRoleById(id: number, role: UserRole): Promise<void> {
    return httpClient.put(`/users/${id}/role`, { role });
  }

  delete(id: number): Promise<void> {
    return httpClient.delete(`/users/${id}`);
  }
}

export const userRepository = new UserRepository();
