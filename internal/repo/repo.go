package repo

import (
	"context"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, delta int64) error
	SetBlocked(ctx context.Context, id uuid.UUID, blocked bool) error
}

type TaskRepo interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Task, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error
	AssignExecutor(ctx context.Context, taskID, executorID uuid.UUID) error
}

type BidRepo interface {
	Create(ctx context.Context, bid *domain.Bid) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Bid, error)
	GetByTaskAndExecutor(ctx context.Context, taskID, executorID uuid.UUID) (*domain.Bid, error)
	ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]*domain.Bid, error)
	HasAcceptedBid(ctx context.Context, taskID uuid.UUID) (bool, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BidStatus) error
	RejectOtherBids(ctx context.Context, taskID, acceptedBidID uuid.UUID) error
}

type PaymentRepo interface {
	Create(ctx context.Context, payment *domain.Payment) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, error)
}
