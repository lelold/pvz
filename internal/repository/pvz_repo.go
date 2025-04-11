package repository

import (
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

type PVZRepository interface {
	Create(pvz *model.PVZ) error
	GetAll() ([]model.PVZ, error)
}

type PVZRepo struct {
	db *gorm.DB
}

func NewPVZRepo(db *gorm.DB) *PVZRepo {
	return &PVZRepo{db: db}
}

func (r *PVZRepo) Create(pvz *model.PVZ) error {
	return r.db.Create(pvz).Error
}

func (r *PVZRepo) GetAll() ([]model.PVZ, error) {
	var pvzList []model.PVZ
	err := r.db.Find(&pvzList).Error
	return pvzList, err
}
