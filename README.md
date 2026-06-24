# Finance & Pomodoro Extension

This repository contains a minimal starter for a **finance logger** and **Pomodoro tracker** designed to be used as an extension for chat agents such as Hermes or OpenClaw.

The architecture separates concerns into three parts:

1. **Backend API (Golang)** – Provides REST endpoints for recording transactions and retrieving summaries.
2. **Database (Supabase/PostgreSQL)** – Stores persistent data such as transactions and Pomodoro sessions.
3. **Frontend Dashboard (Next.js)** – Visualises your financial data and focus statistics.

The goal of this project is to give you a clean foundation without bundling any dependencies. You can run `go mod download` and `npm install` locally to fetch libraries when you are ready.

## Project structure

```text
sregep/
├── backend/         # Golang REST API
│   ├── main.go      # HTTP handlers and Supabase integration
│   ├── go.mod       # Go module definitions
│   └── .env.example # Backend environment template
├── frontend/        # Next.js dashboard starter
│   ├── package.json # Frontend dependencies only, no node_modules
│   ├── src/
│   └── .env.example # Frontend environment template
├── docs/            # Additional documentation
├── .gitignore
└── README.md
```

## Requirements

| Requirement | Version | Purpose |
|---|---:|---|
| Go | 1.21+ | Backend API |
| Node.js | 18+ | Frontend app |
| Supabase | - | PostgreSQL database |
| Git | - | Version control |

## Database schema

Run this SQL in Supabase.

```sql
create table public.transactions (
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

create table public.pomodoro_sessions (
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
```

## Backend API

The backend uses Gin for routing and a Supabase client to interact with Supabase.

### Run backend

```bash
cd backend
cp .env.example .env
# fill SUPABASE_URL and SUPABASE_API_KEY
go mod download
go run main.go
```

The backend listens on `:8080` by default.

### API endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/transactions` | Create transaction |
| GET | `/api/transactions?user_id=<uuid>` | List transactions |
| GET | `/api/summary?user_id=<uuid>` | Get finance summary |

## Frontend

The frontend is a minimal Next.js dashboard starter.

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

## Hermes/OpenClaw integration idea

The first version can expose the Go API through a REST bridge. Later, you can wrap the same finance and Pomodoro operations as MCP tools.

Example AI command:

```text
Catat pengeluaran 25 ribu buat kopi.
```

Expected tool/action payload:

```json
{
  "user_id": "<uuid>",
  "type": "expense",
  "amount": 25000,
  "currency": "IDR",
  "category": "food",
  "note": "kopi",
  "source": "ai"
}
```

## Notes

This is a starter project. Before production use, add authentication, authorization, request validation, RLS policies, rate limiting, and better error handling.
