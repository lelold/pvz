package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Role     string    `gorm:"not null"`
}

type PVZ struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	RegistrationDate time.Time `gorm:"not null"`
	City             string    `gorm:"not null"`
}

type Reception struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateTime time.Time `gorm:"not null"`
	PVZID    uuid.UUID `gorm:"type:uuid;not null"`
	Status   string    `gorm:"not null"`
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateTime    time.Time `gorm:"not null"`
	Type        string    `gorm:"not null"`
	ReceptionID uuid.UUID `gorm:"type:uuid;not null"`
}
