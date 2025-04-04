package interfaces

import (
    "context"
    "mydonate/internal/models"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	
	
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Verify(ctx context.Context, email string, verificationCode string) error
	Update(ctx context.Context, user *models.User) error
	
}