package main

import (
	"log/slog"
	"os"

	"github.com/0xzcarlet/sregep/mcp-server/internal/backend"
	"github.com/0xzcarlet/sregep/mcp-server/internal/config"
	"github.com/0xzcarlet/sregep/mcp-server/internal/mcp"
	"github.com/0xzcarlet/sregep/mcp-server/internal/tools"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	backendClient := backend.NewClient(cfg.APIBaseURL)
	registry := tools.NewRegistry(backendClient, cfg.DefaultUserID)
	server := mcp.NewServer("sregep-mcp", "0.1.0", registry, logger)

	server.Run(os.Stdin, os.Stdout)
}
