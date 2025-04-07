package utils

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

// GenerateRandomString генерирует случайную строку заданной длины.
func GenerateRandomString(length int) string {
    b := make([]byte, length)
    _, err := rand.Read(b)
    if err != nil {
        return ""
    }
    return base64.StdEncoding.EncodeToString(b)
}

// HashPassword хеширует пароль, используя bcrypt.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(bytes), nil
}