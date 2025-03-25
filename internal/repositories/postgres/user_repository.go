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
	var user models.User //просто переменная для хранения результата
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
	if err == gorm.ErrRecordNotFound {
		return nil,nil
	}
	return nil,err
	}
	return &user,nil
}

func (r *UserRepository)  GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository)  UpdateVerificationCode(email string, code string) error {
	//тут не нужена переменная, апдейт возращает СТАТУС
	result := r.db.Model(&models.User{}).Where("email = ?", email).Update("index", code)

	if result.Error != nil{
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound //  Пользователь с таким email не найден
	}
	return nil
}

func (r *UserRepository)  MarkVerified(email string) error {
	result := r.db.Model(&models.User{}).Where("email = ?", email).Update("verified", true)
	if result.Error != nil{
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound //  Пользователь с таким email не найден
	}
	return nil
}






