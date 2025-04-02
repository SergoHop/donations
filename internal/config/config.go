package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadConfig() {
    err := godotenv.Load()
    if err != nil && os.Getenv("ENVIRONMENT") != "production" {
        log.Println("Error loading .env file")
    }
}