'use client';

import { useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { NodeListItem } from '../../internal/domain/model/node';
import type { NodeRecordMap, NodeTrailMap } from '../../internal/domain/service/http/record';
import { MAP_DEFAULT_CENTER, MAP_DEFAULT_ZOOM, TILE_URL, TILE_ATTRIBUTION } from '../../internal/config/constants';
import { formatHeartRate, formatBodyTemp, formatTime } from '../../internal/utils/format';

// Distinct palette — one colour per climber, cycles if more than 8
const TRAIL_PALETTE = [
  '#3b82f6', // blue
  '#f59e0b', // amber
  '#8b5cf6', // violet
  '#ec4899', // pink
  '#06b6d4', // cyan
  '#84cc16', // lime
  '#f97316', // orange
  '#14b8a6', // teal
];

function trailColor(nodeId: number): string {
  return TRAIL_PALETTE[nodeId % TRAIL_PALETTE.length];
}

function createIcon(selected: boolean, color: string) {
  const border = selected ? '#1e3a5f' : darken(color);
  const bg = selected ? '#ef4444' : color;
  return L.divIcon({
    className: '',
    html: `<div style="
      width:28px;height:28px;border-radius:50%;
      background:${bg};border:3px solid ${border};
      box-shadow:0 2px 6px rgba(0,0,0,0.35);
      display:flex;align-items:center;justify-content:center;
    ">
      <div style="width:8px;height:8px;border-radius:50%;background:white;"></div>
    </div>`,
    iconSize: [28, 28],
    iconAnchor: [14, 14],
    popupAnchor: [0, -16],
  });
}

function darken(hex: string): string {
  // Darken a hex colour by ~30% for the marker border
  const n = parseInt(hex.slice(1), 16);
  const r = Math.max(0, ((n >> 16) & 0xff) - 70);
  const g = Math.max(0, ((n >> 8) & 0xff) - 70);
  const b = Math.max(0, (n & 0xff) - 70);
  return `#${((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')}`;
}

function FlyTo({ lat, lon }: { lat: number; lon: number }) {
  const map = useMap();
  useEffect(() => {
    map.flyTo([lat, lon], Math.max(map.getZoom(), 12), { duration: 1 });
  }, [lat, lon, map]);
  return null;
}

interface Props {
  nodes: NodeListItem[];
  recordMap: NodeRecordMap;
  trailMap?: NodeTrailMap;
  selectedNodeId: number | null;
  onSelectNode: (id: number) => void;
}

export default function MapView({ nodes, recordMap, trailMap, selectedNodeId, onSelectNode }: Props) {
  const selectedRecord = selectedNodeId != null ? recordMap.get(selectedNodeId) : undefined;

  return (
    <MapContainer
      center={MAP_DEFAULT_CENTER}
      zoom={MAP_DEFAULT_ZOOM}
      className="h-full w-full"
      zoomControl
    >
      <TileLayer url={TILE_URL} attribution={TILE_ATTRIBUTION} />

      {selectedRecord?.latitude != null && selectedRecord?.longitude != null && (
        <FlyTo lat={selectedRecord.latitude} lon={selectedRecord.longitude} />
      )}

      {/* Trail polylines — rendered before markers so they appear underneath */}
      {nodes.map((node) => {
        const trail = trailMap?.get(node.id);
        if (!trail || trail.length < 2) return null;
        const color = trailColor(node.id);
        return (
          <Polyline
            key={`trail-${node.id}`}
            positions={trail}
            pathOptions={{ color, weight: 3, opacity: 0.65 }}
          />
        );
      })}

      {/* Node markers */}
      {nodes.map((node) => {
        const rec = recordMap.get(node.id);
        if (rec?.latitude == null || rec?.longitude == null) return null;
        const isSelected = node.id === selectedNodeId;
        const color = trailColor(node.id);

        return (
          <Marker
            key={node.id}
            position={[rec.latitude, rec.longitude]}
            icon={createIcon(isSelected, color)}
            eventHandlers={{ click: () => onSelectNode(node.id) }}
          >
            <Popup>
              <div className="text-sm">
                <p className="font-semibold">{node.name}</p>
                <p className="text-slate-500 text-xs">{node.mid}</p>
                <hr className="my-1" />
                <p>HR: {formatHeartRate(rec.heart_rate)}</p>
                <p>Temp: {formatBodyTemp(rec.body_temperature)}</p>
                <p className="text-xs text-slate-400 mt-1">{formatTime(rec.time)}</p>
              </div>
            </Popup>
          </Marker>
        );
      })}
    </MapContainer>
  );
}
