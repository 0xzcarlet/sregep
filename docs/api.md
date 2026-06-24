# API Reference

This document describes the REST API endpoints provided by the Go backend. All endpoints are relative to the server root, for example `http://localhost:8080`.

## Authentication

At this stage, the starter does not enforce user authentication. Before production use, add JWT/API-key validation and Supabase Row Level Security policies.

## Health check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

## Create transaction

```http
POST /api/transactions
```

Create an income or expense transaction.

### Request body

```json
{
  "user_id": "00000000-0000-0000-0000-000000000000",
  "type": "expense",
  "amount": 25000,
  "currency": "IDR",
  "category": "food",
  "note": "kopi",
  "source": "api"
}
```

Required fields:

- `user_id`
- `type`
- `amount`
- `category`

`type` must be `income` or `expense`.

## List transactions

```http
GET /api/transactions?user_id=<uuid>
```

Returns transactions from Supabase ordered by `occurred_at` descending.

## Get summary

```http
GET /api/summary?user_id=<uuid>
```

Response:

```json
{
  "success": true,
  "data": {
    "total_income": 3500000,
    "total_expense": 1250000,
    "balance": 2250000
  }
}
```

## Start Pomodoro

```http
POST /api/pomodoro/start
```

Request body:

```json
{
  "user_id": "00000000-0000-0000-0000-000000000000",
  "task_name": "Build MCP server",
  "duration_minutes": 25
}
```

The backend sets `status` to `running` by default.

## Stop Pomodoro

```http
POST /api/pomodoro/stop
```

Request body:

```json
{
  "user_id": "00000000-0000-0000-0000-000000000000",
  "session_id": "00000000-0000-0000-0000-000000000000"
}
```

The backend updates the matching session to `completed` and fills `ended_at`.

## Current Pomodoro

```http
GET /api/pomodoro/current?user_id=<uuid>
```

Returns the latest running Pomodoro session for the user.

## Next endpoints to add

- `PATCH /api/transactions/:id`
- `DELETE /api/transactions/:id`
- `GET /api/pomodoro/stats`
- `POST /api/pomodoro/cancel`
