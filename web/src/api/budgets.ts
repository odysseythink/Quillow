import client from './client';
import { JsonApiListResponse, JsonApiResponse } from '../types/api';
import { Budget, BudgetLimit } from '../types/models';

export const getBudgets = (page = 1, limit = 50) =>
  client.get<JsonApiListResponse<Budget>>('/budgets', { params: { page, limit } });

export const getBudget = (id: string) =>
  client.get<JsonApiResponse<Budget>>(`/budgets/${id}`);

export const createBudget = (data: Partial<Budget>) =>
  client.post<JsonApiResponse<Budget>>('/budgets', data);

export const updateBudget = (id: string, data: Partial<Budget>) =>
  client.put<JsonApiResponse<Budget>>(`/budgets/${id}`, data);

export const deleteBudget = (id: string) =>
  client.delete(`/budgets/${id}`);

export const getBudgetLimits = (budgetId: string) =>
  client.get<JsonApiListResponse<BudgetLimit>>(`/budgets/${budgetId}/limits`);

export const createBudgetLimit = (budgetId: string, data: Partial<BudgetLimit>) =>
  client.post<JsonApiResponse<BudgetLimit>>(`/budgets/${budgetId}/limits`, data);
