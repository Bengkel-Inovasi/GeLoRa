import type { Record } from '@/internal/domain/model/record';

export interface RecordFilters {
  session_id?: number;
  user_id?: number;
  node_id?: number;
  start_time?: string;
  end_time?: string;
}

export interface IRecordRepository {
  list(filters?: RecordFilters): Promise<{ data: Record[] }>;
}
