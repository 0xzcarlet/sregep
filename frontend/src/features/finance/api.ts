import { env } from '@/lib/env';
import type { ApiResponse, FinanceSummary, Transaction } from './types';

async function request<T>(path: string): Promise<T> {
  const response = await fetch(env.apiBaseUrl + path);
  const payload = (await response.json()) as ApiResponse<T>;

  if (!response.ok || !payload.success || payload.data === undefined) {
    throw new Error(payload.error || 'Request failed');
  }

  return payload.data;
}

export function getFinanceSummary(userId: string) {
  return request<FinanceSummary>('/api/summary?user_id=' + encodeURIComponent(userId));
}

export function getTransactions(userId: string) {
  return request<Transaction[]>('/api/transactions?user_id=' + encodeURIComponent(userId));
}
