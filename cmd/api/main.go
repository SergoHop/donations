package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mydonate/internal/handlers"
	_ "mydonate/internal/interfaces"
	"mydonate/internal/repositories"
	"mydonate/internal/models"
	"mydonate/internal/services"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database configuration
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&models.User{})

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepository)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Gin setup
	router := gin.Default()

	// Routes
	router.POST("/register", userHandler.CreateUserHandler)
	router.POST("/login", userHandler.LoginHandler)
	router.GET("/users/:id", userHandler.GetUserHandler)
	router.PUT("/users/:id", userHandler.UpdateUserHandler)
	router.DELETE("/users/:id", userHandler.DeleteUserHandler)
	router.GET("/verify", userHandler.VerifyHandler) // Добавляем маршрут для подтверждения email

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Starting server on :%s", port)
	router.Run(":" + port)
}