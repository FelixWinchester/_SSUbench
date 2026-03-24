package service

import (
	"context"
	"fmt"
	"time"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/google/uuid"
)

type bidService struct {
	bidRepo  repo.BidRepo
	taskRepo repo.TaskRepo
}

func NewBidService(bidRepo repo.BidRepo, taskRepo repo.TaskRepo) BidService {
	return &bidService{bidRepo: bidRepo, taskRepo: taskRepo}
}

func (s *bidService) Create(ctx context.Context, input CreateBidInput) (*domain.Bid, error) {
	// Задача должна быть опубликована
	task, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if err != nil {
		return nil, fmt.Errorf("bidService.Create get task: %w", err)
	}

	if task.Status != domain.TaskStatusPublished {
		return nil, domain.ErrTaskNotPublished
	}

	// Проверяем что исполнитель ещё не откликался
	_, err = s.bidRepo.GetByTaskAndExecutor(ctx, input.TaskID, input.ExecutorID)
	if err == nil {
		return nil, domain.ErrBidAlreadyExists
	}
	if err != domain.ErrBidNotFound {
		return nil, fmt.Errorf("bidService.Create check existing: %w", err)
	}

	now := time.Now()
	bid := &domain.Bid{
		ID:         uuid.New(),
		TaskID:     input.TaskID,
		ExecutorID: input.ExecutorID,
		Comment:    input.Comment,
		Status:     domain.BidStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err = s.bidRepo.Create(ctx, bid); err != nil {
		return nil, fmt.Errorf("bidService.Create: %w", err)
	}
	return bid, nil
}

func (s *bidService) ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]*domain.Bid, error) {
	bids, err := s.bidRepo.ListByTask(ctx, taskID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("bidService.ListByTask: %w", err)
	}
	return bids, nil
}

func (s *bidService) Accept(ctx context.Context, customerID, taskID, bidID uuid.UUID) error {
	// Проверяем задачу
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("bidService.Accept get task: %w", err)
	}

	if task.CustomerID != customerID {
		return domain.ErrForbidden
	}

	if task.Status != domain.TaskStatusPublished {
		return domain.ErrTaskNotPublished
	}

	// Проверяем что нет уже принятой заявки
	hasAccepted, err := s.bidRepo.HasAcceptedBid(ctx, taskID)
	if err != nil {
		return fmt.Errorf("bidService.Accept check accepted: %w", err)
	}
	if hasAccepted {
		return domain.ErrBidAlreadyAccepted
	}

	// Получаем заявку
	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("bidService.Accept get bid: %w", err)
	}

	if bid.TaskID != taskID {
		return domain.ErrForbidden
	}

	// Принимаем заявку
	if err = s.bidRepo.UpdateStatus(ctx, bidID, domain.BidStatusAccepted); err != nil {
		return fmt.Errorf("bidService.Accept update bid: %w", err)
	}

	// Отклоняем остальные заявки
	if err = s.bidRepo.RejectOtherBids(ctx, taskID, bidID); err != nil {
		return fmt.Errorf("bidService.Accept reject others: %w", err)
	}

	// Назначаем исполнителя на задачу
	if err = s.taskRepo.AssignExecutor(ctx, taskID, bid.ExecutorID); err != nil {
		return fmt.Errorf("bidService.Accept assign executor: %w", err)
	}

	return nil
}
