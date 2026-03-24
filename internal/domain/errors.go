package domain

import "errors"

var (
	// Auth
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")

	// Access
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")

	// Not found
	ErrUserNotFound    = errors.New("user not found")
	ErrTaskNotFound    = errors.New("task not found")
	ErrBidNotFound     = errors.New("bid not found")
	ErrPaymentNotFound = errors.New("payment not found")

	// User
	ErrUserBlocked = errors.New("user is blocked")

	// Task
	ErrTaskInvalidStatus    = errors.New("invalid task status transition")
	ErrTaskAlreadyCompleted = errors.New("task is already completed")
	ErrTaskCancelled        = errors.New("task is cancelled")
	ErrTaskNotPublished     = errors.New("task is not published")
	ErrTaskNotInProgress    = errors.New("task is not in progress")

	// Bid
	ErrBidAlreadyExists   = errors.New("bid already exists")
	ErrBidAlreadyAccepted = errors.New("task already has an accepted bid")

	// Payment
	ErrInsufficientBalance = errors.New("insufficient balance")
)
