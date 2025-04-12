package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Role     string    `gorm:"not null"`
}
