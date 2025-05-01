package models

import (
	"gorm.io/gorm"
)

type User struct {
	Name  string `gorm:"not null" json:"name"`
	Email string `gorm:"uniqueIndex;not null" json:"email"`
	gorm.Model

}

type Wallet struct {
	UserID  uint    `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance float64 `gorm:"not null;default:0" json:"balance"`
	gorm.Model

}
