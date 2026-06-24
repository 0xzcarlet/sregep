create table if not exists public.transactions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null,
  type text not null check (type in ('income','expense')),
  amount numeric not null check (amount > 0),
  currency text not null default 'IDR',
  category text not null,
  note text,
  source text not null default 'api',
  occurred_at timestamptz not null default now(),
  created_at timestamptz not null default now()
);

create table if not exists public.pomodoro_sessions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null,
  task_name text,
  status text not null check (status in ('running','paused','completed','cancelled')),
  duration_minutes integer not null default 25,
  started_at timestamptz not null default now(),
  ended_at timestamptz,
  created_at timestamptz not null default now()
);

alter table public.transactions enable row level security;
alter table public.pomodoro_sessions enable row level security;
