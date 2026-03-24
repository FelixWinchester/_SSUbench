package domain

import (
	"time"

	"github.com/google/uuid"
)

type BidStatus string

const (
	BidStatusPending  BidStatus = "pending"
	BidStatusAccepted BidStatus = "accepted"
	BidStatusRejected BidStatus = "rejected"
)

type Bid struct {
	ID         uuid.UUID `db:"id"`
	TaskID     uuid.UUID `db:"task_id"`
	ExecutorID uuid.UUID `db:"executor_id"`
	Comment    string    `db:"comment"`
	Status     BidStatus `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
