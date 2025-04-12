package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateTime    time.Time `gorm:"not null"`
	Type        string    `gorm:"not null"`
	ReceptionID uuid.UUID `gorm:"type:uuid;not null"`
}
