package supabase

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/0xzcarlet/sregep/backend/internal/domain"
)

type PomodoroRepository struct {
	client *Client
}

func NewPomodoroRepository(client *Client) *PomodoroRepository {
	return &PomodoroRepository{client: client}
}

func (r *PomodoroRepository) Start(ctx context.Context, input domain.PomodoroSession) ([]domain.PomodoroSession, error) {
	var result []domain.PomodoroSession
	err := r.client.Post(ctx, "/rest/v1/pomodoro_sessions", input, &result)
	return result, err
}

func (r *PomodoroRepository) Stop(ctx context.Context, input domain.StopPomodoroInput) ([]domain.PomodoroSession, error) {
	now := time.Now()
	payload := map[string]any{
		"status":   domain.PomodoroCompleted,
		"ended_at": now,
	}
	path := fmt.Sprintf("/rest/v1/pomodoro_sessions?id=eq.%s&user_id=eq.%s", url.QueryEscape(input.SessionID), url.QueryEscape(input.UserID))
	var result []domain.PomodoroSession
	err := r.client.Patch(ctx, path, payload, &result)
	return result, err
}

func (r *PomodoroRepository) Current(ctx context.Context, userID string) ([]domain.PomodoroSession, error) {
	path := fmt.Sprintf("/rest/v1/pomodoro_sessions?user_id=eq.%s&status=eq.running&select=*&order=started_at.desc&limit=1", url.QueryEscape(userID))
	var result []domain.PomodoroSession
	err := r.client.Get(ctx, path, &result)
	return result, err
}
