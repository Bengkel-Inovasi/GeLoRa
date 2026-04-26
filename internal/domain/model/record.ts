export interface Record {
  id: number;
  session_id: number | null;
  time: string;
  heart_rate: number | null;
  body_temperature: number | null;
  latitude: number | null;
  longitude: number | null;
}
