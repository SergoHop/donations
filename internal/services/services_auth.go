package services 

import (
    "context"
	"fmt"
	"log"

	"mydonate/internal/models"       
	"mydonate/internal/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}


func NewUserService(userRepo repositories.UserRepository) UserService {
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

    err := s.userRepo.Create(user)
    if err != nil{
        return fmt.Errorf("фаил ошиька: %w", err)
    }
    return nil
}