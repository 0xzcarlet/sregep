package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/0xzcarlet/sregep/backend/internal/domain"
)

type FinanceService interface {
	CreateTransaction(ctx context.Context, input domain.CreateTransactionInput) ([]domain.Transaction, error)
	ListTransactions(ctx context.Context, userID string) ([]domain.Transaction, error)
	Summary(ctx context.Context, userID string) (domain.FinanceSummary, error)
}

type PomodoroService interface {
	Start(ctx context.Context, input domain.StartPomodoroInput) ([]domain.PomodoroSession, error)
	Stop(ctx context.Context, input domain.StopPomodoroInput) ([]domain.PomodoroSession, error)
	Current(ctx context.Context, userID string) ([]domain.PomodoroSession, error)
}

type Handler struct {
	finance  FinanceService
	pomodoro PomodoroService
	logger   *slog.Logger
}

func NewHandler(finance FinanceService, pomodoro PomodoroService, logger *slog.Logger) *Handler {
	return &Handler{finance: finance, pomodoro: pomodoro, logger: logger}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"success": false, "error": message})
}

func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
