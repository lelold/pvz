package model

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateTime time.Time `gorm:"not null"`
	PVZID    uuid.UUID `gorm:"type:uuid;not null"`
	Status   string    `gorm:"not null"`
}
