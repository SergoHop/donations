package models

type User struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	Username         string `json:"user"`
	Password         string `json:"password"`
	Email            string `json:"email" gorm:"uniqueIndex"`
	VerificationCode string `json:"verification_code" gorm:"size:255"` // <-- Важно: размер
	Verified         bool   `json:"verified"`
	Salt             string `json:"salt"`
}

type RegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"pasword"`
	User     string `json:"user"`
}