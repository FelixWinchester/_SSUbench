package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusDraft      TaskStatus = "draft"
	TaskStatusPublished  TaskStatus = "published"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type Task struct {
	ID          uuid.UUID  `db:"id"`
	CustomerID  uuid.UUID  `db:"customer_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Reward      int64      `db:"reward"`
	Status      TaskStatus `db:"status"`
	ExecutorID  *uuid.UUID `db:"executor_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
