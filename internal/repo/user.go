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

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, balance, is_blocked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Balance,
		user.IsBlocked,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user domain.User
	if err := pgxscan.Get(ctx, r.db, &user, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByID: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user domain.User
	if err := pgxscan.Get(ctx, r.db, &user, query, email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	return &user, nil
}

func (r *userRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	var users []*domain.User
	if err := pgxscan.Select(ctx, r.db, &users, query, limit, offset); err != nil {
		return nil, fmt.Errorf("userRepo.List: %w", err)
	}
	return users, nil
}

func (r *userRepo) UpdateBalance(ctx context.Context, id uuid.UUID, delta int64) error {
	query := `
		UPDATE users
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, delta, id)
	if err != nil {
		return fmt.Errorf("userRepo.UpdateBalance: %w", err)
	}
	return nil
}

func (r *userRepo) SetBlocked(ctx context.Context, id uuid.UUID, blocked bool) error {
	query := `
		UPDATE users
		SET is_blocked = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, blocked, id)
	if err != nil {
		return fmt.Errorf("userRepo.SetBlocked: %w", err)
	}
	return nil
}
