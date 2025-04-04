// internal/database/database.go
package database

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "mydonate/internal/config" // Импортируем config
    "mydonate/internal/models"
)

// DB определяет глобальную переменную для доступа к базе данных.
var DB *gorm.DB

// InitDB инициализирует подключение к базе данных.
func InitDB(cfg config.Config) { // <-- Принимаем cfg
    var err error
    //dsn := os.Getenv("DATABASE_URL") // Больше не нужно

    //if dsn == "" {
    //    log.Fatal("DATABASE_URL is not set in the environment")
    //    return // Добавьте return, чтобы предотвратить дальнейшее выполнение, если переменная не задана
    //}
    dsn := cfg.DatabaseURL // <-- Используем cfg.DatabaseURL

    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    fmt.Println("Connected to database")

    // AutoMigrate выполняет автоматическую миграцию схемы базы данных.
    err = DB.AutoMigrate(&models.User{}) // Перечислите здесь все ваши модели
    if err != nil {
        log.Fatalf("Failed to auto migrate database: %v", err)
    }

    fmt.Println("Database migrated successfully")
}

// GetDB возвращает экземпляр базы данных.
func GetDB() *gorm.DB {
    return DB
}