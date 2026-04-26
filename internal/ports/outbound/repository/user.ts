import type { User, UserListItem, UserRole } from '@/internal/domain/model/user';
import type { PaginationMeta } from '@/internal/domain/model/error';

export interface IUserRepository {
  me(): Promise<User>;
  list(params?: { page?: number; limit?: number }): Promise<{ data: UserListItem[]; pagination: PaginationMeta }>;
  getById(id: number): Promise<User>;
  patchMe(payload: { name?: string; bio?: string }): Promise<void>;
  patchById(id: number, payload: { name?: string; bio?: string }): Promise<void>;
  updatePasswordMe(oldPassword: string, newPassword: string): Promise<void>;
  updatePasswordById(id: number, newPassword: string): Promise<void>;
  updateRoleById(id: number, role: UserRole): Promise<void>;
  delete(id: number): Promise<void>;
}
