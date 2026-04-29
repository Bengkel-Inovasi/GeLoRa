export interface Alert {
  id: number;
  user_id: number | null;
  message: string;
  acknowledged_at: string | null;
  created_at: string;
}
