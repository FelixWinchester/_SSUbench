package tests

import (
	"context"
	"testing"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/mocks"
	"github.com/FelixWinchester/ssubench/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Block_Success(t *testing.T) {
	userRepo := new(mocks.UserRepo)
	svc := service.NewUserService(userRepo)

	adminID := uuid.New()
	userID := uuid.New()

	user := &domain.User{ID: userID}

	userRepo.On("GetByID", context.Background(), userID).Return(user, nil)
	userRepo.On("SetBlocked", context.Background(), userID, true).Return(nil)

	err := svc.Block(context.Background(), adminID, userID)
	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_Block_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepo)
	svc := service.NewUserService(userRepo)

	adminID := uuid.New()
	userID := uuid.New()

	userRepo.On("GetByID", context.Background(), userID).
		Return(nil, domain.ErrUserNotFound)

	err := svc.Block(context.Background(), adminID, userID)
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestPaymentService_ListByUser_Success(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepo)
	svc := service.NewPaymentService(paymentRepo)

	userID := uuid.New()
	expected := []*domain.Payment{
		{ID: uuid.New(), Amount: 100},
		{ID: uuid.New(), Amount: 200},
	}

	paymentRepo.On("ListByUser", context.Background(), userID, 20, 0).
		Return(expected, nil)

	payments, err := svc.ListByUser(context.Background(), userID, 20, 0)
	require.NoError(t, err)
	assert.Len(t, payments, 2)
}
