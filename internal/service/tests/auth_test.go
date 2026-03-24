package tests

import (
	"context"
	"testing"
	"time"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/mocks"
	"github.com/FelixWinchester/ssubench/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := new(mocks.UserRepo)
	svc := service.NewAuthService(userRepo, "secret", time.Hour)

	input := service.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Role:     domain.RoleCustomer,
	}

	userRepo.On("GetByEmail", context.Background(), input.Email).
		Return(nil, domain.ErrUserNotFound)
	userRepo.On("Create", context.Background(), mock.AnythingOfType("*domain.User")).
		Return(nil)

	user, err := svc.Register(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, input.Email, user.Email)
	assert.Equal(t, input.Role, user.Role)
	assert.NotEqual(t, input.Password, user.PasswordHash)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	userRepo := new(mocks.UserRepo)
	svc := service.NewAuthService(userRepo, "secret", time.Hour)

	input := service.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Role:     domain.RoleCustomer,
	}

	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: input.Email,
	}

	userRepo.On("GetByEmail", context.Background(), input.Email).
		Return(existingUser, nil)

	_, err := svc.Register(context.Background(), input)
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userRepo := new(mocks.UserRepo)
	svc := service.NewAuthService(userRepo, "secret", time.Hour)

	input := service.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	existingUser := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: "$2a$12$invalidhash",
		IsBlocked:    false,
	}

	userRepo.On("GetByEmail", context.Background(), input.Email).
		Return(existingUser, nil)

	_, err := svc.Login(context.Background(), input)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestAuthService_Login_BlockedUser(t *testing.T) {
	userRepo := new(mocks.UserRepo)
	svc := service.NewAuthService(userRepo, "secret", time.Hour)

	input := service.LoginInput{
		Email:    "blocked@example.com",
		Password: "password123",
	}

	blockedUser := &domain.User{
		ID:        uuid.New(),
		Email:     input.Email,
		IsBlocked: true,
	}

	userRepo.On("GetByEmail", context.Background(), input.Email).
		Return(blockedUser, nil)

	_, err := svc.Login(context.Background(), input)
	assert.ErrorIs(t, err, domain.ErrUserBlocked)
}
