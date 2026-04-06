import client from './client';
import { JsonApiListResponse, JsonApiResponse } from '../types/api';
import { Account } from '../types/models';

export const getAccounts = (page = 1, limit = 50, type?: string) =>
  client.get<JsonApiListResponse<Account>>('/accounts', { params: { page, limit, type } });

export const getAccount = (id: string) =>
  client.get<JsonApiResponse<Account>>(`/accounts/${id}`);

export const createAccount = (data: Partial<Account> & { type: string }) =>
  client.post<JsonApiResponse<Account>>('/accounts', data);

export const updateAccount = (id: string, data: Partial<Account>) =>
  client.put<JsonApiResponse<Account>>(`/accounts/${id}`, data);

export const deleteAccount = (id: string) =>
  client.delete(`/accounts/${id}`);
