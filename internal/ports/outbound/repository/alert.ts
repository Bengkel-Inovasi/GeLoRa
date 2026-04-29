import type { Alert } from '@/internal/domain/model/alert';

export interface IAlertRepository {
  create(message: string): Promise<{ id: number }>;
  list(params?: { unacknowledged?: boolean }): Promise<{ data: Alert[] }>;
  acknowledge(id: number): Promise<void>;
}
