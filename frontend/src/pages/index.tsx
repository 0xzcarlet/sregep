import { useEffect, useState } from 'react';

type Summary = {
  total_income: number;
  total_expense: number;
  balance: number;
};

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
const demoUserId = process.env.NEXT_PUBLIC_DEMO_USER_ID || '';

export default function Home() {
  const [summary, setSummary] = useState<Summary | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchSummary() {
      if (!demoUserId) {
        setError('Set NEXT_PUBLIC_DEMO_USER_ID in frontend/.env.local first.');
        return;
      }

      try {
        const res = await fetch(`${apiBaseUrl}/api/summary?user_id=${demoUserId}`);
        const json = await res.json();

        if (!json.success) {
          setError(json.error || 'Failed to fetch summary.');
          return;
        }

        setSummary(json.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error.');
      }
    }

    fetchSummary();
  }, []);

  return (
    <main className="mx-auto max-w-5xl px-6 py-10">
      <div className="mb-8">
        <p className="text-sm font-medium text-slate-500">Sregep</p>
        <h1 className="text-3xl font-bold">Finance & Pomodoro Dashboard</h1>
        <p className="mt-2 text-slate-600">
          Starter dashboard for finance logging and focus tracking.
        </p>
      </div>

      {error && (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-red-700">
          {error}
        </div>
      )}

      {!summary && !error && <p>Loading summary...</p>}

      {summary && (
        <div className="grid gap-4 md:grid-cols-3">
          <SummaryCard label="Income" value={summary.total_income} />
          <SummaryCard label="Expense" value={summary.total_expense} />
          <SummaryCard label="Balance" value={summary.balance} />
        </div>
      )}
    </main>
  );
}

function SummaryCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-2xl border bg-white p-5 shadow-sm">
      <p className="text-sm text-slate-500">{label}</p>
      <p className="mt-2 text-2xl font-bold">Rp{value.toLocaleString('id-ID')}</p>
    </div>
  );
}
