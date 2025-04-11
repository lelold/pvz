package repository

import (
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

type ReceptionRepo struct {
	db *gorm.DB
}

func NewReceptionRepo(db *gorm.DB) *ReceptionRepo {
	return &ReceptionRepo{db: db}
}

func (r *ReceptionRepo) CreateReception(reception *model.Reception) error {
	return r.db.Create(reception).Error
}

func (r *ReceptionRepo) HasOpenReception(pvzID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Reception{}).
		Where("pvz_id = ? AND status = ?", pvzID, "close").
		Count(&count).Error
	return count > 0, err
}
