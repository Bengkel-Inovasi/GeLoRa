export interface ApiError {
  error: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
}
