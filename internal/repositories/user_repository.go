package repositories

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	
	"mydonate/internal/models"
	"mydonate/internal/interfaces"
)


type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository { // <-- Принимаем db
    return &userRepository{db: db}
}
// Create создает нового пользователя в базе данных.
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}
// получает пользователя по имаил
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("ошибка с имаилом: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("ошибка с айди: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("ошибка с обновление пользователя: %w", result.Error)
	}
	return nil
}

func (r *userRepository) Verify(ctx context.Context, email string, verificationCode string) error {
	return nil
}






