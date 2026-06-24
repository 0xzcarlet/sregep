package service

import (
	"context"
	"errors"
	"strings"

	"github.com/0xzcarlet/sregep/backend/internal/domain"
)

type FinanceRepository interface {
	Create(ctx context.Context, input domain.CreateTransactionInput) ([]domain.Transaction, error)
	ListByUserID(ctx context.Context, userID string) ([]domain.Transaction, error)
}

type FinanceService struct {
	repo FinanceRepository
}

func NewFinanceService(repo FinanceRepository) *FinanceService {
	return &FinanceService{repo: repo}
}

func (s *FinanceService) CreateTransaction(ctx context.Context, input domain.CreateTransactionInput) ([]domain.Transaction, error) {
	input.UserID = strings.TrimSpace(input.UserID)
	input.Category = strings.TrimSpace(input.Category)

	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if input.Type != domain.TransactionIncome && input.Type != domain.TransactionExpense {
		return nil, errors.New("type must be income or expense")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if input.Category == "" {
		return nil, errors.New("category is required")
	}
	if input.Currency == "" {
		input.Currency = "IDR"
	}
	if input.Source == "" {
		input.Source = "api"
	}

	return s.repo.Create(ctx, input)
}

func (s *FinanceService) ListTransactions(ctx context.Context, userID string) ([]domain.Transaction, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.ListByUserID(ctx, userID)
}

func (s *FinanceService) Summary(ctx context.Context, userID string) (domain.FinanceSummary, error) {
	transactions, err := s.ListTransactions(ctx, userID)
	if err != nil {
		return domain.FinanceSummary{}, err
	}

	var summary domain.FinanceSummary
	for _, trx := range transactions {
		switch trx.Type {
		case domain.TransactionIncome:
			summary.TotalIncome += trx.Amount
		case domain.TransactionExpense:
			summary.TotalExpense += trx.Amount
		}
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}
