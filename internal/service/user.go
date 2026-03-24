package service

import (
	"context"
	"fmt"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/google/uuid"
)

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userService.GetByID: %w", err)
	}
	return user, nil
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("userService.List: %w", err)
	}
	return users, nil
}

func (s *userService) Block(ctx context.Context, adminID, userID uuid.UUID) error {
	// Проверяем что пользователь существует
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("userService.Block: %w", err)
	}

	if err = s.userRepo.SetBlocked(ctx, userID, true); err != nil {
		return fmt.Errorf("userService.Block: %w", err)
	}
	return nil
}

func (s *userService) Unblock(ctx context.Context, adminID, userID uuid.UUID) error {
	// Проверяем что пользователь существует
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("userService.Unblock: %w", err)
	}

	if err = s.userRepo.SetBlocked(ctx, userID, false); err != nil {
		return fmt.Errorf("userService.Unblock: %w", err)
	}
	return nil
}
