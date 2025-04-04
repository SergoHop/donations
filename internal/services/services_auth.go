package services 

import (
    "context"
	"fmt"
	"log"

    "github.com/go-playground/validator/v10"
	"mydonate/internal/models"       
	"mydonate/internal/interfaces"
)

type userService struct {
	userRepository interfaces.UserRepository
    validator      *validator.Validate 
}


func NewUserService(userRepository interfaces.UserRepository) interfaces.UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
    if user == nil{
        return fmt.Errorf("юзер наможа быть нил")
    }
    if user.Email == "" {
		return fmt.Errorf("шо то не то")
	}

    err := s.userRepository.Create(ctx, user)
    if err != nil{
        return fmt.Errorf("фаил ошиька: %w", err)
    }
    return nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    if email == ""{
        return nil, fmt.Errorf("шо то не то")
    }
    user, err := s.userRepository.GetByEmail(ctx, email)
    if err != nil{
        return nil, fmt.Errorf("нема такого имаила %w", err)
    }
    return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
    if id == 0{
        return nil, fmt.Errorf("гдэ айди")
    }
    user, err := s.userRepository.GetByID(ctx, id)
    if err != nil{
        return nil, fmt.Errorf("а шо с айди %w", err)
    }
    return user, nil
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	//Валидация данных пользователя перед обновлением
	err := s.validator.Struct(user)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return s.userRepository.Update(ctx, user)
}

func (s *userService) Verify(ctx context.Context, email string, verificationCode string) error {
	// Получаем пользователя по email
	user, err := s.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Проверяем, совпадает ли код верификации
	if user.VerificationCode != verificationCode {
		return fmt.Errorf("invalid verification code")
	}

	if user.Verified {
		return fmt.Errorf("user already verified")
	}

	// Обновляем статус пользователя как верифицированного
	user.Verified = true
	err = s.Update(ctx, user) // Используем метод Update для сохранения изменений
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("User %s verified successfully", email)
	return nil
}