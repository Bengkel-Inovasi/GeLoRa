import type { IRecordRepository, RecordFilters } from '../../../ports/outbound/repository/record';
import type { Record } from '../../../domain/model/record';
import { httpClient } from '../client/http';

class RecordRepository implements IRecordRepository {
  list(filters?: RecordFilters): Promise<{ data: Record[] }> {
    const q = new URLSearchParams();
    if (filters?.session_id !== undefined) q.set('session_id', String(filters.session_id));
    if (filters?.user_id !== undefined) q.set('user_id', String(filters.user_id));
    if (filters?.node_id !== undefined) q.set('node_id', String(filters.node_id));
    if (filters?.start_time) q.set('start_time', filters.start_time);
    if (filters?.end_time) q.set('end_time', filters.end_time);
    const qs = q.toString() ? `?${q}` : '';
    return httpClient.get(`/records${qs}`);
  }
}

export const recordRepository = new RecordRepository();
