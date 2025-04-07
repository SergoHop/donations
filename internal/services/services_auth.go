package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"mydonate/internal/interfaces"
	"mydonate/internal/models"
	"mydonate/internal/utils"

	"github.com/jordan-wright/email"
	"gorm.io/gorm"

	"net/smtp"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(userRepository interfaces.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	existingUser, err := s.userRepository.GetByEmail(ctx, user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Hash the password
	salt := utils.GenerateRandomString(16)
	hashedPassword, err := utils.HashPassword(user.Password + salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	verificationCode := utils.GenerateRandomString(48)

	user.Password = hashedPassword
	user.Salt = salt
	user.VerificationCode = verificationCode
	user.Verified = false

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = SendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		return fmt.Errorf("Failed to send verification email")
	}

	return nil
}

func (s *UserService) Get(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *UserService) VerifyEmail(ctx context.Context, email string, code string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.VerificationCode != code {
		return fmt.Errorf("invalid verification code")
	}

	user.Verified = true
	return s.userRepository.Update(ctx, user)
}

func SendVerificationEmail(emailAddress string, verificationCode string) error {
	from := os.Getenv("SMTP_FROM")
	log.Printf("SMTP_FROM: %s", from)
	password := os.Getenv("SMTP_PASSWORD")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")

	e := email.NewEmail()
	e.From = from
	e.To = []string{emailAddress}
	e.Subject = "Подтверждение регистрации"
	e.HTML = []byte(fmt.Sprintf(`<h1>Здравствуйте!</h1><p>Для подтверждения регистрации, перейдите по ссылке: <a href="http://localhost:8080/verify?email=%s&code=%s">Подтвердить</a></p>`, emailAddress, verificationCode))

	auth := smtp.PlainAuth("", from, password, smtpServer)
	addr := smtpServer + ":" + smtpPort
	err := e.Send(addr, auth)
	if err != nil {
		return err
	}

	log.Println("Verification email sent to " + emailAddress)
	return nil
}