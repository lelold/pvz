package model

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	RegistrationDate time.Time `gorm:"not null"`
	City             string    `gorm:"not null"`
}
