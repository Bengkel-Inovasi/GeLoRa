export interface Session {
  id: number;
  user_id: number | null;
  node_id: number | null;
  started_at: string;
  ended_at: string | null;
}
