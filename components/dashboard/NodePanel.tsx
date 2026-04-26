'use client';

import type { NodeListItem } from '../../internal/domain/model/node';
import type { Record } from '../../internal/domain/model/record';
import Card from '../ui/Card';
import Badge from '../ui/Badge';
import { formatHeartRate, formatBodyTemp, formatCoords, formatTime } from '../../internal/utils/format';

interface Props {
  node: NodeListItem | null;
  record: Record | undefined;
}

function StatRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between py-2 border-b border-slate-100 last:border-0">
      <span className="text-xs text-slate-500 uppercase tracking-wide">{label}</span>
      <span className="text-sm font-semibold text-slate-800">{value}</span>
    </div>
  );
}

export default function NodePanel({ node, record }: Props) {
  if (!node) {
    return (
      <Card className="h-full">
        <div className="flex h-full items-center justify-center text-slate-400 text-sm">
          Select a device on the map to view data
        </div>
      </Card>
    );
  }

  return (
    <Card className="h-full overflow-y-auto">
      <div className="flex items-start justify-between mb-4">
        <div>
          <p className="font-semibold text-slate-800">{node.name}</p>
          <p className="text-xs text-slate-400 font-mono mt-0.5">{node.mid}</p>
        </div>
        <Badge
          label={node.is_validated ? 'Validated' : 'Unvalidated'}
          variant={node.is_validated ? 'green' : 'red'}
        />
      </div>

      {!record ? (
        <p className="text-sm text-slate-400 italic">No recent data received</p>
      ) : (
        <>
          <div className="mb-3">
            <StatRow label="Heart Rate" value={formatHeartRate(record.heart_rate)} />
            <StatRow label="Body Temp" value={formatBodyTemp(record.body_temperature)} />
            <StatRow label="Coordinates" value={formatCoords(record.latitude, record.longitude)} />
            <StatRow label="Last Update" value={formatTime(record.time)} />
          </div>

          {record.latitude != null && record.longitude != null && (
            <a
              href={`https://www.openstreetmap.org/?mlat=${record.latitude}&mlon=${record.longitude}&zoom=15`}
              target="_blank"
              rel="noopener noreferrer"
              className="mt-2 inline-block text-xs text-emerald-600 underline"
            >
              Open in OSM
            </a>
          )}
        </>
      )}
    </Card>
  );
}
