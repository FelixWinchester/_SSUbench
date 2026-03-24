package mocks

import (
	"context"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type BidRepo struct {
	mock.Mock
}

func (m *BidRepo) Create(ctx context.Context, bid *domain.Bid) error {
	args := m.Called(ctx, bid)
	return args.Error(0)
}

func (m *BidRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bid, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bid), args.Error(1)
}

func (m *BidRepo) GetByTaskAndExecutor(ctx context.Context, taskID, executorID uuid.UUID) (*domain.Bid, error) {
	args := m.Called(ctx, taskID, executorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bid), args.Error(1)
}

func (m *BidRepo) ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]*domain.Bid, error) {
	args := m.Called(ctx, taskID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Bid), args.Error(1)
}

func (m *BidRepo) HasAcceptedBid(ctx context.Context, taskID uuid.UUID) (bool, error) {
	args := m.Called(ctx, taskID)
	return args.Bool(0), args.Error(1)
}

func (m *BidRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BidStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *BidRepo) RejectOtherBids(ctx context.Context, taskID, acceptedBidID uuid.UUID) error {
	args := m.Called(ctx, taskID, acceptedBidID)
	return args.Error(0)
}
