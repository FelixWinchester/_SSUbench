package mocks

import (
	"context"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type TaskRepo struct {
	mock.Mock
}

func (m *TaskRepo) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskRepo) List(ctx context.Context, limit, offset int) ([]*domain.Task, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *TaskRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *TaskRepo) AssignExecutor(ctx context.Context, taskID, executorID uuid.UUID) error {
	args := m.Called(ctx, taskID, executorID)
	return args.Error(0)
}
