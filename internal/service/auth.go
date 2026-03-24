package service

import (
	"context"
	"fmt"
	"time"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  repo.UserRepo
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(userRepo repo.UserRepo, jwtSecret string, jwtTTL time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	// Проверяем что email не занят
	_, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("authService.Register: %w", err)
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("authService.Register bcrypt: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         input.Role,
		Balance:      0,
		IsBlocked:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("authService.Register: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if user.IsBlocked {
		return "", domain.ErrUserBlocked
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", fmt.Errorf("authService.Login: %w", err)
	}

	return token, nil
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *authService) generateToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID.String(),
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
