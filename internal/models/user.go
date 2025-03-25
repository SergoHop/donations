package models

import (
	"errors"
	"gorm.io/gorm"
)

var ErrNoRecord = errors.New("записи нема)")

type User struct{
	gorm.Model
	Username string    `gorm:"not null" json:"user"`
	Password string          `gorm:"not null" json:"pasword"`
	Email string             `gorm:"uniqueIndex" json:"enail"`
	VerificationCode string  `gorm:"index"`
	Verified bool            `gorm:"default:false"`
}

