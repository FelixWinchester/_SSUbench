package service

import (
	"context"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input LoginInput) (string, error)
}

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	Block(ctx context.Context, adminID, userID uuid.UUID) error
	Unblock(ctx context.Context, adminID, userID uuid.UUID) error
}

type TaskService interface {
	Create(ctx context.Context, input CreateTaskInput) (*domain.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Task, error)
	Publish(ctx context.Context, customerID, taskID uuid.UUID) error
	Cancel(ctx context.Context, customerID, taskID uuid.UUID) error
	MarkCompleted(ctx context.Context, executorID, taskID uuid.UUID) error
	Confirm(ctx context.Context, customerID, taskID uuid.UUID) error
}

type BidService interface {
	Create(ctx context.Context, input CreateBidInput) (*domain.Bid, error)
	ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]*domain.Bid, error)
	Accept(ctx context.Context, customerID, taskID, bidID uuid.UUID) error
}

type PaymentService interface {
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, error)
}

// Input structs

type RegisterInput struct {
	Email    string      `json:"email"    validate:"required,email"`
	Password string      `json:"password" validate:"required,min=6"`
	Role     domain.Role `json:"role"     validate:"required,oneof=customer executor"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateTaskInput struct {
	CustomerID  uuid.UUID `validate:"required"`
	Title       string    `json:"title"       validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"required,min=10"`
	Reward      int64     `json:"reward"      validate:"required,min=1"`
}

type CreateBidInput struct {
	TaskID     uuid.UUID `validate:"required"`
	ExecutorID uuid.UUID `validate:"required"`
	Comment    string    `json:"comment" validate:"max=1000"`
}
