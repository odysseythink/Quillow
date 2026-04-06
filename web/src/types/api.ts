export interface JsonApiResource<T> {
  type: string;
  id: string;
  attributes: T;
}

export interface JsonApiResponse<T> {
  data: JsonApiResource<T>;
}

export interface JsonApiListResponse<T> {
  data: JsonApiResource<T>[];
  meta: {
    pagination: PaginationMeta;
  };
}

export interface PaginationMeta {
  total: number;
  count: number;
  per_page: number;
  current_page: number;
  total_pages: number;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
}

export interface ErrorResponse {
  message: string;
  exception: number;
}
