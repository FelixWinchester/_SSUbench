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

type bidRepo struct {
	db *pgxpool.Pool
}

func NewBidRepo(db *pgxpool.Pool) BidRepo {
	return &bidRepo{db: db}
}

func (r *bidRepo) Create(ctx context.Context, bid *domain.Bid) error {
	query := `
		INSERT INTO bids (id, task_id, executor_id, comment, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		bid.ID,
		bid.TaskID,
		bid.ExecutorID,
		bid.Comment,
		bid.Status,
		bid.CreatedAt,
		bid.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("bidRepo.Create: %w", err)
	}
	return nil
}

func (r *bidRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bid, error) {
	query := `SELECT * FROM bids WHERE id = $1`

	var bid domain.Bid
	if err := pgxscan.Get(ctx, r.db, &bid, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBidNotFound
		}
		return nil, fmt.Errorf("bidRepo.GetByID: %w", err)
	}
	return &bid, nil
}

func (r *bidRepo) GetByTaskAndExecutor(ctx context.Context, taskID, executorID uuid.UUID) (*domain.Bid, error) {
	query := `SELECT * FROM bids WHERE task_id = $1 AND executor_id = $2`

	var bid domain.Bid
	if err := pgxscan.Get(ctx, r.db, &bid, query, taskID, executorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBidNotFound
		}
		return nil, fmt.Errorf("bidRepo.GetByTaskAndExecutor: %w", err)
	}
	return &bid, nil
}

func (r *bidRepo) ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]*domain.Bid, error) {
	query := `SELECT * FROM bids WHERE task_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var bids []*domain.Bid
	if err := pgxscan.Select(ctx, r.db, &bids, query, taskID, limit, offset); err != nil {
		return nil, fmt.Errorf("bidRepo.ListByTask: %w", err)
	}
	return bids, nil
}

func (r *bidRepo) HasAcceptedBid(ctx context.Context, taskID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM bids WHERE task_id = $1 AND status = 'accepted')`

	var exists bool
	if err := r.db.QueryRow(ctx, query, taskID).Scan(&exists); err != nil {
		return false, fmt.Errorf("bidRepo.HasAcceptedBid: %w", err)
	}
	return exists, nil
}

func (r *bidRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BidStatus) error {
	query := `
		UPDATE bids
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("bidRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *bidRepo) RejectOtherBids(ctx context.Context, taskID, acceptedBidID uuid.UUID) error {
	query := `
		UPDATE bids
		SET status = 'rejected', updated_at = NOW()
		WHERE task_id = $1 AND id != $2 AND status = 'pending'
	`
	_, err := r.db.Exec(ctx, query, taskID, acceptedBidID)
	if err != nil {
		return fmt.Errorf("bidRepo.RejectOtherBids: %w", err)
	}
	return nil
}
