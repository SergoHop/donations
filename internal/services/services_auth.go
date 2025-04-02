package services 

import (
    "context"
	"fmt"
	"log"

	"mydonate/internal/models"       
	"mydonate/internal/repositories"
)


type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	UpdateVerificationCode(ctx context.Context, email string, code string) error
	MarkVerified(ctx context.Context, email string) error
}

type userService struct {
	userRepo repositories.UserService
}


func NewUserService(userRepo repositories.UserService) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
    if user == nil{
        return fmt.Errorf("юзер наможа быть нил")
    }
    if user.Email == "" {
		return fmt.Errorf("шо то не то")
	}

    err := s.userRepo.Create(ctx, user)
    if err != nil{
        return fmt.Errorf("фаил ошиька: %w", err)
    }
    return nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    if email == ""{
        return nil, fmt.Errorf("шо то не то")
    }
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil{
        return nil, fmt.Errorf("нема такого имаила %w", err)
    }
    return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
    if id == 0{
        return nil, fmt.Errorf("гдэ айди")
    }
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil{
        return nil, fmt.Errorf("а шо с айди %w", err)
    }
    return user, nil
}

func (s *userService) UpdateVerificationCode(ctx context.Context, email string, code string) error{
    if email == ""{
        return fmt.Errorf("нема такого имаила")
    }
    if code == ""{
        return fmt.Errorf("а шо с кодом")
    }
    err := s.userRepo.UpdateVerificationCode(ctx, email, code)
    if err != nil{
        log.Printf("Ошибка при обновлении кода верификации: %v", err)
		return fmt.Errorf("ошибка при обновлении кода верификации: %w", err)
    }
    return nil
}

func (s *userService) MarkVerified(ctx context.Context, email string) error{
    if email == ""{
        return fmt.Errorf("нема такого имаила")
    }
    err := s.userRepo.MarkVerified(ctx, email)
	if err != nil {
		log.Printf("Ошибка при подтверждении пользователя: %v", err)
		return fmt.Errorf("ошибка при подтверждении пользователя: %w", err)
    }
    return nil
}