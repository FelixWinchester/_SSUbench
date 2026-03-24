package service

import (
	"context"
	"fmt"
	"time"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// txTaskService расширяет taskService возможностью выполнять транзакции
type txTaskService struct {
	*taskService
	db *pgxpool.Pool
}

func NewTxTaskService(
	db *pgxpool.Pool,
	taskRepo repo.TaskRepo,
	userRepo repo.UserRepo,
	bidRepo repo.BidRepo,
	paymentRepo repo.PaymentRepo,
) TaskService {
	base := &taskService{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		bidRepo:     bidRepo,
		paymentRepo: paymentRepo,
	}
	return &txTaskService{taskService: base, db: db}
}

func (s *txTaskService) Confirm(ctx context.Context, customerID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm: %w", err)
	}

	if task.CustomerID != customerID {
		return domain.ErrForbidden
	}

	if task.Status != domain.TaskStatusCompleted {
		return domain.ErrTaskInvalidStatus
	}

	customer, err := s.userRepo.GetByID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm get customer: %w", err)
	}

	if customer.Balance < task.Reward {
		return domain.ErrInsufficientBalance
	}

	if task.ExecutorID == nil {
		return domain.ErrForbidden
	}

	executorID := *task.ExecutorID

	// Атомарная транзакция
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Списываем у заказчика
	_, err = tx.Exec(ctx,
		`UPDATE users SET balance = balance - $1, updated_at = NOW() WHERE id = $2 AND balance >= $1`,
		task.Reward, customerID,
	)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm debit: %w", err)
	}

	// 2. Начисляем исполнителю
	_, err = tx.Exec(ctx,
		`UPDATE users SET balance = balance + $1, updated_at = NOW() WHERE id = $2`,
		task.Reward, executorID,
	)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm credit: %w", err)
	}

	// 3. Записываем платёж
	_, err = tx.Exec(ctx,
		`INSERT INTO payments (id, task_id, from_user_id, to_user_id, amount, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), task.ID, customerID, executorID, task.Reward, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm insert payment: %w", err)
	}

	// 4. Обновляем статус задачи
	_, err = tx.Exec(ctx,
		`UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2`,
		domain.TaskStatusCompleted, taskID,
	)
	if err != nil {
		return fmt.Errorf("txTaskService.Confirm update task: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("txTaskService.Confirm commit: %w", err)
	}

	return nil
}
