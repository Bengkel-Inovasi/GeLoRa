'use client';

import { useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { NodeListItem } from '../../internal/domain/model/node';
import type { NodeRecordMap } from '../../internal/domain/service/http/record';
import { MAP_DEFAULT_CENTER, MAP_DEFAULT_ZOOM, TILE_URL, TILE_ATTRIBUTION } from '../../internal/config/constants';
import { formatHeartRate, formatBodyTemp, formatCoords, formatTime } from '../../internal/utils/format';

function createIcon(selected: boolean) {
  const color = selected ? '#ef4444' : '#059669';
  const border = selected ? '#991b1b' : '#065f46';
  return L.divIcon({
    className: '',
    html: `<div style="
      width:28px;height:28px;border-radius:50%;
      background:${color};border:3px solid ${border};
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
  selectedNodeId: number | null;
  onSelectNode: (id: number) => void;
}

export default function MapView({ nodes, recordMap, selectedNodeId, onSelectNode }: Props) {
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

      {nodes.map((node) => {
        const rec = recordMap.get(node.id);
        if (rec?.latitude == null || rec?.longitude == null) return null;
        const isSelected = node.id === selectedNodeId;

        return (
          <Marker
            key={node.id}
            position={[rec.latitude, rec.longitude]}
            icon={createIcon(isSelected)}
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
