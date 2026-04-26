export type UserRole = 'super' | 'admin' | 'client';

export interface User {
  id: number;
  name: string;
  username: string;
  role: UserRole;
  bio: string;
  created_at: string;
  updated_at: string | null;
}

export interface UserListItem {
  id: number;
  name: string;
  username: string;
  role: UserRole;
}
