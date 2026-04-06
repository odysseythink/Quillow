import client from './client';
import { JsonApiListResponse, JsonApiResponse } from '../types/api';
import { TransactionGroup } from '../types/models';

export const getTransactions = (page = 1, limit = 50, type?: string, start?: string, end?: string) =>
  client.get<JsonApiListResponse<TransactionGroup>>('/transactions', { params: { page, limit, type, start, end } });

export const getTransaction = (id: string) =>
  client.get<JsonApiResponse<TransactionGroup>>(`/transactions/${id}`);

export const createTransaction = (data: any) =>
  client.post<JsonApiResponse<TransactionGroup>>('/transactions', data);

export const deleteTransaction = (id: string) =>
  client.delete(`/transactions/${id}`);

export const searchTransactions = (query: string, page = 1, limit = 50) =>
  client.get<JsonApiListResponse<TransactionGroup>>('/search/transactions', { params: { query, page, limit } });

export const getSummary = (start?: string, end?: string) =>
  client.get('/summary/basic', { params: { start, end } });
