package domain

import "time"

type TransactionType string

const (
	TransactionIncome  TransactionType = "income"
	TransactionExpense TransactionType = "expense"
)

type Transaction struct {
	ID         string          `json:"id,omitempty"`
	UserID     string          `json:"user_id"`
	Type       TransactionType `json:"type"`
	Amount     float64         `json:"amount"`
	Currency   string          `json:"currency"`
	Category   string          `json:"category"`
	Note       string          `json:"note,omitempty"`
	Source     string          `json:"source"`
	OccurredAt *time.Time      `json:"occurred_at,omitempty"`
	CreatedAt  string          `json:"created_at,omitempty"`
}

type CreateTransactionInput struct {
	UserID     string          `json:"user_id"`
	Type       TransactionType `json:"type"`
	Amount     float64         `json:"amount"`
	Currency   string          `json:"currency"`
	Category   string          `json:"category"`
	Note       string          `json:"note"`
	Source     string          `json:"source"`
	OccurredAt *time.Time      `json:"occurred_at,omitempty"`
}

type FinanceSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type PomodoroStatus string

const (
	PomodoroRunning   PomodoroStatus = "running"
	PomodoroPaused    PomodoroStatus = "paused"
	PomodoroCompleted PomodoroStatus = "completed"
	PomodoroCancelled PomodoroStatus = "cancelled"
)

type PomodoroSession struct {
	ID              string         `json:"id,omitempty"`
	UserID          string         `json:"user_id"`
	TaskName        string         `json:"task_name,omitempty"`
	Status          PomodoroStatus `json:"status"`
	DurationMinutes int            `json:"duration_minutes"`
	StartedAt       *time.Time     `json:"started_at,omitempty"`
	EndedAt         *time.Time     `json:"ended_at,omitempty"`
	CreatedAt       string         `json:"created_at,omitempty"`
}

type StartPomodoroInput struct {
	UserID          string `json:"user_id"`
	TaskName        string `json:"task_name"`
	DurationMinutes int    `json:"duration_minutes"`
}

type StopPomodoroInput struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}
