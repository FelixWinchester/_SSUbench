package service

import (
	"context"
	"fmt"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/google/uuid"
)

type paymentService struct {
	paymentRepo repo.PaymentRepo
}

func NewPaymentService(paymentRepo repo.PaymentRepo) PaymentService {
	return &paymentService{paymentRepo: paymentRepo}
}

func (s *paymentService) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, error) {
	payments, err := s.paymentRepo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("paymentService.ListByUser: %w", err)
	}
	return payments, nil
}
