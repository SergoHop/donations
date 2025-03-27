package repositories

import (
    "context"
    "mydonate/internal/models"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	UpdateVerificationCode(ctx context.Context, email string, code string) error
	MarkVerified(ctx context.Context, email string) error
}