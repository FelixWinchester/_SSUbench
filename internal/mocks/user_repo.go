package mocks

import (
	"context"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepo struct {
	mock.Mock
}

func (m *UserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *UserRepo) UpdateBalance(ctx context.Context, id uuid.UUID, delta int64) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func (m *UserRepo) SetBlocked(ctx context.Context, id uuid.UUID, blocked bool) error {
	args := m.Called(ctx, id, blocked)
	return args.Error(0)
}
