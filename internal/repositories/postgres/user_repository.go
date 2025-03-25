package postgres

import (
	"gorm.io/gorm"
	"mydonate/internal/models"
	
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
// Create создает нового пользователя в базе данных.
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
// получает пользователя по имаил
func (r *UserRepository)  GetByEmail(email string) (*models.User, error) {
	
}

func (r *UserRepository)  GetByID(id uint) (*models.User, error) {

}

func (r *UserRepository)  UpdateVerificationCode(email string, code string) error {

}

func (r *UserRepository)  MarkVerified(email string) error {

}






