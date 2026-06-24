# MCP Server for Hermes

This project includes a Go-based MCP stdio server under `mcp-server/`.

The MCP server does not write directly to Supabase. It calls the Go backend API, and the backend writes to Supabase.

```text
Hermes / MCP client
  -> mcp-server
  -> backend API
  -> Supabase
```

## Available MCP tools

| Tool | Purpose |
|---|---|
| `finance_add_transaction` | Record income or expense |
| `finance_list_transactions` | List user transactions |
| `finance_summary` | Get income, expense, balance |
| `pomodoro_start` | Start focus session |
| `pomodoro_stop` | Stop session by ID |
| `pomodoro_current` | Get running session |

## 1. Run backend first

```bash
cd backend
cp .env.example .env
# fill SUPABASE_URL and SUPABASE_API_KEY
go mod download
go run main.go
```

Backend default URL:

```text
http://localhost:8080
```

## 2. Run MCP server manually

```bash
cd mcp-server
cp .env.example .env
export SREGEP_API_BASE_URL=http://localhost:8080
export SREGEP_DEFAULT_USER_ID=your-user-uuid

go run main.go
```

The MCP server talks over stdio, so it will wait for JSON-RPC messages from Hermes or another MCP client.

## 3. Build MCP binary

```bash
cd mcp-server
go build -o sregep-mcp .
```

Then run:

```bash
SREGEP_API_BASE_URL=http://localhost:8080 \
SREGEP_DEFAULT_USER_ID=your-user-uuid \
./sregep-mcp
```

## 4. Example Hermes MCP config

Exact config depends on your Hermes setup. The core idea is to register this command as a stdio MCP server.

```json
{
  "mcpServers": {
    "sregep": {
      "command": "/absolute/path/to/sregep/mcp-server/sregep-mcp",
      "args": [],
      "env": {
        "SREGEP_API_BASE_URL": "http://localhost:8080",
        "SREGEP_DEFAULT_USER_ID": "your-user-uuid"
      }
    }
  }
}
```

If Hermes expects `go run`, use this during development:

```json
{
  "mcpServers": {
    "sregep": {
      "command": "go",
      "args": ["run", "/absolute/path/to/sregep/mcp-server/main.go"],
      "env": {
        "SREGEP_API_BASE_URL": "http://localhost:8080",
        "SREGEP_DEFAULT_USER_ID": "your-user-uuid"
      }
    }
  }
}
```

## 5. Example prompts

```text
Catat pengeluaran 25 ribu untuk kopi, kategori food.
```

Expected tool call:

```json
{
  "name": "finance_add_transaction",
  "arguments": {
    "type": "expense",
    "amount": 25000,
    "category": "food",
    "note": "kopi"
  }
}
```

```text
Mulai pomodoro 25 menit buat ngerjain Go API.
```

Expected tool call:

```json
{
  "name": "pomodoro_start",
  "arguments": {
    "task_name": "ngerjain Go API",
    "duration_minutes": 25
  }
}
```

## Notes

- Set `SREGEP_DEFAULT_USER_ID` to avoid passing `user_id` in every prompt.
- Keep the backend API running before starting the MCP server.
- For production, run backend and MCP server via `systemd`, Docker, or process manager.
- Do not expose Supabase service role keys to frontend or public clients.
