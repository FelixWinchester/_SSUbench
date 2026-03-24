package mocks

import (
	"context"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type PaymentRepo struct {
	mock.Mock
}

func (m *PaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *PaymentRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Payment), args.Error(1)
}
