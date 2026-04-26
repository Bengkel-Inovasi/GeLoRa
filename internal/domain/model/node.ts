export interface Node {
  id: number;
  mid: string;
  name: string;
  description: string;
  is_validated: boolean;
  validated_by: number | null;
  created_at: string;
  updated_at: string | null;
}

export interface NodeListItem {
  id: number;
  mid: string;
  name: string;
  is_validated: boolean;
}
