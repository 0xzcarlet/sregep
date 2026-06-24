# Finance & Pomodoro Extension

This repository contains a minimal starter for a **finance logger** and **Pomodoro tracker** designed to be used as an extension for chat agents such as Hermes or OpenClaw.

The architecture separates concerns into four parts:

1. **Backend API (Golang)** – Provides REST endpoints for finance and Pomodoro actions.
2. **MCP Server (Golang)** – Exposes backend actions as MCP tools over stdio.
3. **Database (Supabase/PostgreSQL)** – Stores transactions and Pomodoro sessions.
4. **Frontend Dashboard (Next.js)** – Visualises your financial data and focus statistics.

The goal of this project is to give you a clean foundation without bundling any dependencies. You can run `go mod download` and `npm install` locally to fetch libraries when you are ready.

## Project structure

```text
sregep/
├── backend/            # Golang REST API
│   ├── main.go
│   ├── go.mod
│   └── .env.example
├── mcp-server/         # Golang MCP stdio server for Hermes/OpenClaw
│   ├── main.go
│   ├── go.mod
│   └── .env.example
├── frontend/           # Next.js dashboard starter
│   ├── package.json
│   ├── src/
│   └── .env.example
├── docs/
│   ├── api.md
│   ├── mcp-hermes.md
│   └── schema.sql
├── .gitignore
└── README.md
```

## Requirements

| Requirement | Version | Purpose |
|---|---:|---|
| Go | 1.21+ | Backend and MCP server |
| Node.js | 18+ | Frontend app |
| Supabase | - | PostgreSQL database |
| Git | - | Version control |

## Database schema

Run the SQL in `docs/schema.sql` inside Supabase.

## Backend API

The backend uses Gin for routing and Supabase REST API for persistence.

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
| GET | `/health` | Health check |
| POST | `/api/transactions` | Create transaction |
| GET | `/api/transactions?user_id=<uuid>` | List transactions |
| GET | `/api/summary?user_id=<uuid>` | Get finance summary |
| POST | `/api/pomodoro/start` | Start Pomodoro session |
| POST | `/api/pomodoro/stop` | Stop Pomodoro session |
| GET | `/api/pomodoro/current?user_id=<uuid>` | Get running session |

## MCP server for Hermes

The MCP server lives in `mcp-server/` and talks to the backend API.

### Run MCP server

```bash
cd mcp-server
cp .env.example .env
export SREGEP_API_BASE_URL=http://localhost:8080
export SREGEP_DEFAULT_USER_ID=your-user-uuid

go run main.go
```

### Build MCP binary

```bash
cd mcp-server
go build -o sregep-mcp .
```

### MCP tools

| Tool | Purpose |
|---|---|
| `finance_add_transaction` | Record income or expense |
| `finance_list_transactions` | List transactions |
| `finance_summary` | Get finance summary |
| `pomodoro_start` | Start focus session |
| `pomodoro_stop` | Stop session by ID |
| `pomodoro_current` | Get running session |

Detailed setup is in `docs/mcp-hermes.md`.

## Frontend

The frontend is a minimal Next.js dashboard starter.

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

## Example AI commands

```text
Catat pengeluaran 25 ribu buat kopi kategori food.
```

```text
Mulai pomodoro 25 menit buat ngerjain Go API.
```

```text
Cek summary finance bulan ini.
```

## Notes

This is a starter project. Before production use, add authentication, authorization, request validation, Supabase RLS policies, rate limiting, and better error handling.
