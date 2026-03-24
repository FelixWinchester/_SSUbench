package tests

import (
	"context"
	"testing"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/mocks"
	"github.com/FelixWinchester/ssubench/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBidService_Create_TaskNotPublished(t *testing.T) {
	bidRepo := new(mocks.BidRepo)
	taskRepo := new(mocks.TaskRepo)

	svc := service.NewBidService(bidRepo, taskRepo)

	taskID := uuid.New()
	executorID := uuid.New()

	task := &domain.Task{
		ID:     taskID,
		Status: domain.TaskStatusDraft,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	input := service.CreateBidInput{
		TaskID:     taskID,
		ExecutorID: executorID,
	}

	_, err := svc.Create(context.Background(), input)
	assert.ErrorIs(t, err, domain.ErrTaskNotPublished)
}

func TestBidService_Create_AlreadyExists(t *testing.T) {
	bidRepo := new(mocks.BidRepo)
	taskRepo := new(mocks.TaskRepo)

	svc := service.NewBidService(bidRepo, taskRepo)

	taskID := uuid.New()
	executorID := uuid.New()

	task := &domain.Task{
		ID:     taskID,
		Status: domain.TaskStatusPublished,
	}

	existingBid := &domain.Bid{
		ID:         uuid.New(),
		TaskID:     taskID,
		ExecutorID: executorID,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)
	bidRepo.On("GetByTaskAndExecutor", context.Background(), taskID, executorID).
		Return(existingBid, nil)

	input := service.CreateBidInput{
		TaskID:     taskID,
		ExecutorID: executorID,
	}

	_, err := svc.Create(context.Background(), input)
	assert.ErrorIs(t, err, domain.ErrBidAlreadyExists)
}

func TestBidService_Accept_AlreadyAccepted(t *testing.T) {
	bidRepo := new(mocks.BidRepo)
	taskRepo := new(mocks.TaskRepo)

	svc := service.NewBidService(bidRepo, taskRepo)

	customerID := uuid.New()
	taskID := uuid.New()
	bidID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: customerID,
		Status:     domain.TaskStatusPublished,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)
	bidRepo.On("HasAcceptedBid", context.Background(), taskID).Return(true, nil)

	err := svc.Accept(context.Background(), customerID, taskID, bidID)
	assert.ErrorIs(t, err, domain.ErrBidAlreadyAccepted)
}

func TestBidService_Accept_NotOwner(t *testing.T) {
	bidRepo := new(mocks.BidRepo)
	taskRepo := new(mocks.TaskRepo)

	svc := service.NewBidService(bidRepo, taskRepo)

	customerID := uuid.New()
	otherID := uuid.New()
	taskID := uuid.New()
	bidID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: otherID,
		Status:     domain.TaskStatusPublished,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	err := svc.Accept(context.Background(), customerID, taskID, bidID)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}
