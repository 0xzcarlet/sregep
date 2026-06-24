package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xzcarlet/sregep/backend/internal/config"
	supabaserepo "github.com/0xzcarlet/sregep/backend/internal/repository/supabase"
	"github.com/0xzcarlet/sregep/backend/internal/server"
	"github.com/0xzcarlet/sregep/backend/internal/service"
	"github.com/0xzcarlet/sregep/backend/internal/transport/httpapi"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	supabaseClient := supabaserepo.NewClient(cfg.SupabaseURL, cfg.SupabaseAPIKey)
	financeRepo := supabaserepo.NewFinanceRepository(supabaseClient)
	pomodoroRepo := supabaserepo.NewPomodoroRepository(supabaseClient)

	financeService := service.NewFinanceService(financeRepo)
	pomodoroService := service.NewPomodoroService(pomodoroRepo)

	router := httpapi.NewRouter(httpapi.RouterConfig{
		AllowedOrigin:   cfg.AllowedOrigin,
		FinanceService:  financeService,
		PomodoroService: pomodoroService,
		Logger:          logger,
	})

	httpServer := server.New(cfg.Port, router, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("backend started", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error("backend stopped", "error", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
