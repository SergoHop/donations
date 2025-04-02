package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"mydonate/internal/handlers"
	"mydonate/internal/services"
	"mydonate/internal/repositories"
	"database/sql"
	"os"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
}

defer db.Close()

    userRepository := repositories.NewUserRepository(db) // Создаем репозиторий
    userService := services.NewUserService(userRepository)  // Создаем сервис
    userHandler := handlers.NewUserHandler(userService)   // Создаем handler

	// Создаем экземпляр Gin
	r := gin.Default()

	// Регистрируем handlers для различных маршрутов
	r.POST("/users", userHandler.CreateUserHandler)
	r.GET("/users/:id", userHandler.GetUserByIDHandler)
	r.GET("/users", userHandler.GetUserByEmailHandler)

	// Запускаем HTTP-сервер
	log.Println("сервак запущен на порту: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}