export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: string;
};

export type FinanceSummary = {
  total_income: number;
  total_expense: number;
  balance: number;
};

export type Transaction = {
  id: string;
  user_id: string;
  type: 'income' | 'expense';
  amount: number;
  currency: string;
  category: string;
  note?: string;
  source: string;
  occurred_at?: string;
  created_at?: string;
};
