package repo

import (
	"context"
	"fmt"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type paymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepo(db *pgxpool.Pool) PaymentRepo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, task_id, from_user_id, to_user_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		payment.ID,
		payment.TaskID,
		payment.FromUserID,
		payment.ToUserID,
		payment.Amount,
		payment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("paymentRepo.Create: %w", err)
	}
	return nil
}

func (r *paymentRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, error) {
	query := `
		SELECT * FROM payments
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var payments []*domain.Payment
	if err := pgxscan.Select(ctx, r.db, &payments, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("paymentRepo.ListByUser: %w", err)
	}
	return payments, nil
}
