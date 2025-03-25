package repositories

import "mydonate/internal/models"

type UserRepository interface {
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
    GetByID(id uint) (*models.User, error)
    UpdateVerificationCode(email string, code string) error
	MarkVerified(email string) error
}