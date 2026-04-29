import type { IAlertRepository } from '../../../ports/outbound/repository/alert';
import type { Alert } from '../../../domain/model/alert';
import { httpClient } from '../client/http';

class AlertRepository implements IAlertRepository {
  create(message: string): Promise<{ id: number }> {
    return httpClient.post('/alerts', { message });
  }

  list(params?: { unacknowledged?: boolean }): Promise<{ data: Alert[] }> {
    const unack = params?.unacknowledged ?? true;
    return httpClient.get(`/alerts?unacknowledged=${unack}`);
  }

  acknowledge(id: number): Promise<void> {
    return httpClient.put(`/alerts/${id}/acknowledge`, {});
  }
}

export const alertRepository = new AlertRepository();
