package tests

import (
	"context"
	"testing"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/mocks"
	"github.com/FelixWinchester/ssubench/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_Publish_NotOwner(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	customerID := uuid.New()
	otherID := uuid.New()
	taskID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: otherID,
		Status:     domain.TaskStatusDraft,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	err := svc.Publish(context.Background(), customerID, taskID)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTaskService_Publish_WrongStatus(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	customerID := uuid.New()
	taskID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: customerID,
		Status:     domain.TaskStatusPublished,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	err := svc.Publish(context.Background(), customerID, taskID)
	assert.ErrorIs(t, err, domain.ErrTaskInvalidStatus)
}

func TestTaskService_Cancel_CompletedTask(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	customerID := uuid.New()
	taskID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: customerID,
		Status:     domain.TaskStatusCompleted,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	err := svc.Cancel(context.Background(), customerID, taskID)
	assert.ErrorIs(t, err, domain.ErrTaskAlreadyCompleted)
}

func TestTaskService_MarkCompleted_NotExecutor(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	executorID := uuid.New()
	otherExecutorID := uuid.New()
	taskID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		Status:     domain.TaskStatusInProgress,
		ExecutorID: &otherExecutorID,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	err := svc.MarkCompleted(context.Background(), executorID, taskID)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTaskService_Confirm_InsufficientBalance(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	customerID := uuid.New()
	executorID := uuid.New()
	taskID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: customerID,
		Status:     domain.TaskStatusCompleted,
		Reward:     1000,
		ExecutorID: &executorID,
	}

	customer := &domain.User{
		ID:      customerID,
		Balance: 500,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)
	userRepo.On("GetByID", context.Background(), customerID).Return(customer, nil)

	err := svc.Confirm(context.Background(), customerID, taskID)
	assert.ErrorIs(t, err, domain.ErrInsufficientBalance)
}

func TestTaskService_Confirm_NotOwner(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	customerID := uuid.New()
	otherID := uuid.New()
	taskID := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		CustomerID: otherID,
		Status:     domain.TaskStatusCompleted,
	}

	taskRepo.On("GetByID", context.Background(), taskID).Return(task, nil)

	err := svc.Confirm(context.Background(), customerID, taskID)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTaskService_Create_Success(t *testing.T) {
	taskRepo := new(mocks.TaskRepo)
	userRepo := new(mocks.UserRepo)
	bidRepo := new(mocks.BidRepo)
	paymentRepo := new(mocks.PaymentRepo)

	svc := service.NewTaskService(taskRepo, userRepo, bidRepo, paymentRepo)

	input := service.CreateTaskInput{
		CustomerID:  uuid.New(),
		Title:       "Test Task",
		Description: "Test Description long enough",
		Reward:      100,
	}

	taskRepo.On("Create", context.Background(), mock.AnythingOfType("*domain.Task")).
		Return(nil)

	task, err := svc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.Title, task.Title)
	assert.Equal(t, domain.TaskStatusDraft, task.Status)
}
