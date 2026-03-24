package service

import (
	"context"
	"fmt"
	"time"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/google/uuid"
)

type taskService struct {
	taskRepo    repo.TaskRepo
	userRepo    repo.UserRepo
	bidRepo     repo.BidRepo
	paymentRepo repo.PaymentRepo
}

func NewTaskService(
	taskRepo repo.TaskRepo,
	userRepo repo.UserRepo,
	bidRepo repo.BidRepo,
	paymentRepo repo.PaymentRepo,
) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		bidRepo:     bidRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *taskService) Create(ctx context.Context, input CreateTaskInput) (*domain.Task, error) {
	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New(),
		CustomerID:  input.CustomerID,
		Title:       input.Title,
		Description: input.Description,
		Reward:      input.Reward,
		Status:      domain.TaskStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("taskService.Create: %w", err)
	}
	return task, nil
}

func (s *taskService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("taskService.GetByID: %w", err)
	}
	return task, nil
}

func (s *taskService) List(ctx context.Context, limit, offset int) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("taskService.List: %w", err)
	}
	return tasks, nil
}

func (s *taskService) Publish(ctx context.Context, customerID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("taskService.Publish: %w", err)
	}

	if task.CustomerID != customerID {
		return domain.ErrForbidden
	}

	if task.Status != domain.TaskStatusDraft {
		return domain.ErrTaskInvalidStatus
	}

	return s.taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusPublished)
}

func (s *taskService) Cancel(ctx context.Context, customerID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("taskService.Cancel: %w", err)
	}

	if task.CustomerID != customerID {
		return domain.ErrForbidden
	}

	if task.Status == domain.TaskStatusCompleted {
		return domain.ErrTaskAlreadyCompleted
	}

	if task.Status == domain.TaskStatusCancelled {
		return domain.ErrTaskCancelled
	}

	return s.taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusCancelled)
}

func (s *taskService) MarkCompleted(ctx context.Context, executorID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("taskService.MarkCompleted: %w", err)
	}

	if task.Status != domain.TaskStatusInProgress {
		return domain.ErrTaskNotInProgress
	}

	// Только назначенный исполнитель может пометить задачу выполненной
	if task.ExecutorID == nil || *task.ExecutorID != executorID {
		return domain.ErrForbidden
	}

	return s.taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusCompleted)
}

func (s *taskService) Confirm(ctx context.Context, customerID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("taskService.Confirm: %w", err)
	}

	if task.CustomerID != customerID {
		return domain.ErrForbidden
	}

	if task.Status != domain.TaskStatusCompleted {
		return domain.ErrTaskInvalidStatus
	}

	customer, err := s.userRepo.GetByID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("taskService.Confirm get customer: %w", err)
	}

	if customer.Balance < task.Reward {
		return domain.ErrInsufficientBalance
	}

	if task.ExecutorID == nil {
		return domain.ErrForbidden
	}

	return nil
}
