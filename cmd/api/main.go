package main

import (
	
	"log"
	
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mydonate/internal/config"
	"mydonate/internal/handlers"
	
	"mydonate/internal/models"
	"mydonate/internal/repositories"
	"mydonate/internal/services"
)

func main() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация GORM
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate - создает таблицы на основе структуры моделей
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Инициализация репозитория
	userRepository := repositories.NewUserRepository(db)

	// Инициализация сервиса
	userService := services.NewUserService(userRepository)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)

	// Создание Gin router
	router := gin.Default()

	// Определение маршрутов
	router.POST("/register", userHandler.CreateUserHandler)
	router.GET("/verify", userHandler.VerifyEmailHandler)
	router.GET("/users/:id", userHandler.GetUserByIDHandler)
	router.GET("/users", userHandler.GetUserByEmailHandler)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(router.Run(":" + port))
}