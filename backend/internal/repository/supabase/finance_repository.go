package supabase

import (
	"context"
	"fmt"
	"net/url"

	"github.com/0xzcarlet/sregep/backend/internal/domain"
)

type FinanceRepository struct {
	client *Client
}

func NewFinanceRepository(client *Client) *FinanceRepository {
	return &FinanceRepository{client: client}
}

func (r *FinanceRepository) Create(ctx context.Context, input domain.CreateTransactionInput) ([]domain.Transaction, error) {
	var result []domain.Transaction
	err := r.client.Post(ctx, "/rest/v1/transactions", input, &result)
	return result, err
}

func (r *FinanceRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Transaction, error) {
	path := fmt.Sprintf("/rest/v1/transactions?user_id=eq.%s&select=*&order=occurred_at.desc", url.QueryEscape(userID))
	var result []domain.Transaction
	err := r.client.Get(ctx, path, &result)
	return result, err
}
