"use client";

import { useEffect, useState } from 'react';
import { EmptyState } from '@/components/empty-state';
import { SummaryCard } from '@/components/summary-card';
import { getFinanceSummary, getTransactions } from '@/features/finance/api';
import type { FinanceSummary, Transaction } from '@/features/finance/types';
import { env } from '@/lib/env';
import { formatCurrency } from '@/lib/format';

export function DashboardPage() {
  const [summary, setSummary] = useState<FinanceSummary | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function loadDashboard() {
      if (!env.demoUserId) {
        setError('Set NEXT_PUBLIC_DEMO_USER_ID in frontend/.env.local first.');
        setIsLoading(false);
        return;
      }

      try {
        const [summaryData, transactionData] = await Promise.all([
          getFinanceSummary(env.demoUserId),
          getTransactions(env.demoUserId)
        ]);

        setSummary(summaryData);
        setTransactions(transactionData);
      } catch (caught) {
        setError(caught instanceof Error ? caught.message : 'Unknown error');
      } finally {
        setIsLoading(false);
      }
    }

    loadDashboard();
  }, []);

  return (
    <main className="mx-auto flex min-h-screen max-w-6xl flex-col gap-8 px-6 py-10">
      <header>
        <p className="text-sm font-semibold uppercase tracking-wide text-slate-500">Sregep</p>
        <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-950">Finance & Pomodoro Dashboard</h1>
        <p className="mt-3 max-w-2xl text-slate-600">
          Dashboard starter untuk finance logging dan focus tracking yang bisa dipanggil dari Hermes lewat MCP.
        </p>
      </header>

      {isLoading ? <p className="text-slate-600">Loading dashboard...</p> : null}

      {error ? (
        <div className="rounded-2xl border border-red-200 bg-red-50 p-5 text-sm text-red-700">
          {error}
        </div>
      ) : null}

      {summary ? (
        <section className="grid gap-4 md:grid-cols-3">
          <SummaryCard label="Income" value={formatCurrency(summary.total_income)} />
          <SummaryCard label="Expense" value={formatCurrency(summary.total_expense)} />
          <SummaryCard label="Balance" value={formatCurrency(summary.balance)} />
        </section>
      ) : null}

      <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-950">Recent Transactions</h2>
          <p className="text-sm text-slate-500">{transactions.length} item</p>
        </div>

        {transactions.length === 0 && !isLoading ? (
          <EmptyState title="No transaction yet" description="Catat transaksi pertama dari Hermes atau API backend." />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="text-slate-500">
                <tr>
                  <th className="py-3">Type</th>
                  <th className="py-3">Category</th>
                  <th className="py-3">Note</th>
                  <th className="py-3 text-right">Amount</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {transactions.map((transaction) => (
                  <tr key={transaction.id}>
                    <td className="py-3 capitalize">{transaction.type}</td>
                    <td className="py-3">{transaction.category}</td>
                    <td className="py-3 text-slate-500">{transaction.note || '-'}</td>
                    <td className="py-3 text-right font-medium">{formatCurrency(transaction.amount)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </main>
  );
}
