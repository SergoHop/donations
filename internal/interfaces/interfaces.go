package interfaces

import (
	"context"
	"mydonate/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	VerifyEmail(ctx context.Context, email string, code string) error
}