package httpapi

import (
	"log/slog"
	"net/http"
)

type RouterConfig struct {
	AllowedOrigin   string
	FinanceService  FinanceService
	PomodoroService PomodoroService
	Logger          *slog.Logger
}

func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()
	h := NewHandler(cfg.FinanceService, cfg.PomodoroService, cfg.Logger)

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/transactions", h.Transactions)
	mux.HandleFunc("/api/summary", h.Summary)
	mux.HandleFunc("/api/pomodoro/start", h.StartPomodoro)
	mux.HandleFunc("/api/pomodoro/stop", h.StopPomodoro)
	mux.HandleFunc("/api/pomodoro/current", h.CurrentPomodoro)

	return withCORS(mux, cfg.AllowedOrigin)
}

func withCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
