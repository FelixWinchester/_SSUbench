package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) TaskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, customer_id, executor_id, title, description, reward, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		task.ID,
		task.CustomerID,
		task.ExecutorID,
		task.Title,
		task.Description,
		task.Reward,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("taskRepo.Create: %w", err)
	}
	return nil
}

func (r *taskRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1`

	var task domain.Task
	if err := pgxscan.Get(ctx, r.db, &task, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("taskRepo.GetByID: %w", err)
	}
	return &task, nil
}

func (r *taskRepo) List(ctx context.Context, limit, offset int) ([]*domain.Task, error) {
	query := `SELECT * FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	var tasks []*domain.Task
	if err := pgxscan.Select(ctx, r.db, &tasks, query, limit, offset); err != nil {
		return nil, fmt.Errorf("taskRepo.List: %w", err)
	}
	return tasks, nil
}

func (r *taskRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	query := `
		UPDATE tasks
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("taskRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *taskRepo) AssignExecutor(ctx context.Context, taskID, executorID uuid.UUID) error {
	query := `
		UPDATE tasks
		SET executor_id = $1, status = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, executorID, domain.TaskStatusInProgress, taskID)
	if err != nil {
		return fmt.Errorf("taskRepo.AssignExecutor: %w", err)
	}
	return nil
}
