package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/0xzcarlet/sregep/backend/internal/domain"
)

type PomodoroRepository interface {
	Start(ctx context.Context, input domain.PomodoroSession) ([]domain.PomodoroSession, error)
	Stop(ctx context.Context, input domain.StopPomodoroInput) ([]domain.PomodoroSession, error)
	Current(ctx context.Context, userID string) ([]domain.PomodoroSession, error)
}

type PomodoroService struct {
	repo PomodoroRepository
}

func NewPomodoroService(repo PomodoroRepository) *PomodoroService {
	return &PomodoroService{repo: repo}
}

func (s *PomodoroService) Start(ctx context.Context, input domain.StartPomodoroInput) ([]domain.PomodoroSession, error) {
	input.UserID = strings.TrimSpace(input.UserID)
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if input.DurationMinutes <= 0 {
		input.DurationMinutes = 25
	}
	if strings.TrimSpace(input.TaskName) == "" {
		input.TaskName = "Focus session"
	}

	now := time.Now()
	session := domain.PomodoroSession{
		UserID:          input.UserID,
		TaskName:        input.TaskName,
		Status:          domain.PomodoroRunning,
		DurationMinutes: input.DurationMinutes,
		StartedAt:       &now,
	}
	return s.repo.Start(ctx, session)
}

func (s *PomodoroService) Stop(ctx context.Context, input domain.StopPomodoroInput) ([]domain.PomodoroSession, error) {
	input.UserID = strings.TrimSpace(input.UserID)
	input.SessionID = strings.TrimSpace(input.SessionID)
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if input.SessionID == "" {
		return nil, errors.New("session_id is required")
	}
	return s.repo.Stop(ctx, input)
}

func (s *PomodoroService) Current(ctx context.Context, userID string) ([]domain.PomodoroSession, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.Current(ctx, userID)
}
