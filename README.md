# Sregep

Sregep is a finance logger and Pomodoro tracker that works from a normal dashboard and from AI tools through an MCP server.

## Architecture

```text
Hermes / MCP client
  -> mcp-server
  -> backend API
  -> Supabase

Browser
  -> frontend dashboard
  -> backend API
  -> Supabase
```

## Stack

| Layer | Stack |
|---|---|
| Backend | Go standard library |
| MCP | Go stdio JSON-RPC |
| Database | Supabase PostgreSQL |
| Frontend | Next.js App Router |
| Styling | Tailwind CSS |

## Project structure

```text
sregep/
- backend/cmd/api
- backend/internal/config
- backend/internal/domain
- backend/internal/repository/supabase
- backend/internal/server
- backend/internal/service
- backend/internal/transport/httpapi
- frontend/src/app
- frontend/src/components
- frontend/src/features
- frontend/src/lib
- mcp-server/cmd/server
- mcp-server/internal/backend
- mcp-server/internal/config
- mcp-server/internal/mcp
- mcp-server/internal/tools
- docs
```

## 1. Setup Supabase

Run SQL from:

```bash
cat docs/schema.sql
```

## 2. Run backend

```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

Backend default URL:

```text
http://localhost:8080
```

## 3. Run frontend

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

Frontend default URL:

```text
http://localhost:3000
```

## 4. Run MCP server

```bash
cd mcp-server
cp .env.example .env
export SREGEP_API_BASE_URL=http://localhost:8080
export SREGEP_DEFAULT_USER_ID=your-user-uuid

go run ./cmd/server
```

For Hermes, build binary:

```bash
cd mcp-server
go build -o sregep-mcp ./cmd/server
```

Then register `sregep-mcp` as a stdio MCP server. See `docs/mcp-hermes.md`.

## MCP tools

| Tool | Purpose |
|---|---|
| `finance_add_transaction` | Record income or expense |
| `finance_list_transactions` | List transactions |
| `finance_summary` | Get balance summary |
| `pomodoro_start` | Start focus session |
| `pomodoro_stop` | Stop focus session |
| `pomodoro_current` | Get running session |

## Notes

This starter does not commit dependencies such as `node_modules`. It also keeps local environment files ignored by Git.

Before production use, add real authentication, Supabase RLS policies, request limits, structured logs, and better observability.
