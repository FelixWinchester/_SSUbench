package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID         uuid.UUID `db:"id"`
	TaskID     uuid.UUID `db:"task_id"`
	FromUserID uuid.UUID `db:"from_user_id"`
	ToUserID   uuid.UUID `db:"to_user_id"`
	Amount     int64     `db:"amount"`
	CreatedAt  time.Time `db:"created_at"`
}
