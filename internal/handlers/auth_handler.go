package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"mydonate/internal/interfaces"
	"mydonate/internal/models"
	"net/http"
	"net/smtp"
	"os"
	

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{
    userService interfaces.UserService
}

func NewUserHadler(userService interfaces.UserService) *UserHandler{
    return &UserHandler{userService: userService}
}
// тута мы генерируем случайный код верификации.
func GenerateVerificationCode() (string, error){
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil{
		return "", err

	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func HashPassword(password string) (string, string, error){
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil{
		return "", "", nil
	}
	salt := base64.StdEncoding.EncodeToString(saltBytes)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	hashedPassword := string(hashedPasswordBytes)
	return hashedPassword, salt, nil
}
// отправляет email с кодом верификации.
func SendVerificationEmail(emailAddress string, verificationCode string) error {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")

	e := email.NewEmail()
	e.From = from
	e.To = []string{emailAddress}
	e.Subject = "Подтверждение регистрации"
	e.HTML = []byte(fmt.Sprintf(`<h1>Здравствуйте!</h1><p>Для подтверждения регистрации, перейдите по ссылке: <a href="http://localhost:8080/verify?email=%s&code=%s">Подтвердить</a></p>`, emailAddress, verificationCode))

	auth := smtp.PlainAuth("", from, password, smtpServer)
	err := e.Send(smtpServer+":"+smtpPort, auth)
	if err != nil {
		return err
	}

	log.Println("Verification email sent to " + emailAddress)
	return nil
}

// CreateUserHandler обрабатывает запрос на создание нового пользователя.
// POST /users
func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var user models.User

	// Декодируем JSON из тела запроса в структуру User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Хешируем пароль
	hashedPassword, salt, err := HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Генерируем код верификации
	verificationCode, err := GenerateVerificationCode()
	if err != nil {
		log.Printf("Error generating verification code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification code"})
		return
	}

	// Обновляем данные пользователя
	user.Password = hashedPassword
	user.Salt = salt
	user.VerificationCode = verificationCode
	user.Verified = false // Устанавливаем статус "не подтвержден"

	// Создаем пользователя в базе данных (через сервис)
	err = h.userService.Create(c.Request.Context(), &user)
	if err != nil {
		log.Printf("Error creating user: %v", err)

		// Проверяем, является ли ошибка ошибкой валидации
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Создаем карту ошибок валидации
			errorMap := make(map[string]string)
			for _, fieldError := range validationErrors {
				errorMap[fieldError.Field()] = fieldError.Tag() // Сохраняем имя поля и тег валидации
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMap}) // Отправляем клиенту карту ошибок
		} else {
			// Если это не ошибка валидации, отправляем общую ошибку сервера
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	// Отправляем email с кодом верификации
	err = SendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UserHandler) GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	// Преобразовать id в число (uint) и обработать ошибки
	userID, err := StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}


func (h *UserHandler) GetUserByEmailHandler(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	user, err := h.userService.GetByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) VerifyEmailHandler(c *gin.Context) {
	emailAddress := c.Query("email")
	verificationCode := c.Query("code")

	if emailAddress == "" || verificationCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and code are required"})
		return
	}

	// Получаем пользователя по email
	user, err := h.userService.GetByEmail(c.Request.Context(), emailAddress)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Проверяем код верификации и обновляем статус пользователя
	if user.VerificationCode == verificationCode && !user.Verified {
		user.Verified = true
		err = h.userService.Update(c.Request.Context(), user) // Обновляем пользователя
		if err != nil {
			log.Printf("Error updating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		c.String(http.StatusOK, "Email verified successfully!")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code or user already verified"})
	}
}

// StringToUint преобразует строку в uint.
func StringToUint(s string) (uint, error) {
	var i uint
	_, err := fmt.Sscan(s, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}