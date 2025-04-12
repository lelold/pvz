package repository

import (
	"errors"
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

type ReceptionRepo interface {
	CreateReception(reception *model.Reception) error
	FindLastOpenReceptionByPVZ(pvzID string) (*model.Reception, error)
	CloseReception(reception *model.Reception) error
	HasOpenReception(pvzID string) (bool, error)
}

type receptionRepo struct {
	db *gorm.DB
}

func NewReceptionRepo(db *gorm.DB) *receptionRepo {
	return &receptionRepo{db: db}
}

func (r *receptionRepo) CreateReception(reception *model.Reception) error {
	return r.db.Create(reception).Error
}

func (r *receptionRepo) HasOpenReception(pvzID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Reception{}).
		Where("pvz_id = ? AND status = ?", pvzID, "in_progress").
		Count(&count).Error
	return count > 0, err
}

func (r *receptionRepo) FindLastOpenReceptionByPVZ(pvzID string) (*model.Reception, error) {
	var reception model.Reception
	err := r.db.
		Where("pvz_id = ? AND status = ?", pvzID, "in_progress").
		Order("date_time DESC").
		First(&reception).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("открытая приемка не найдена")
		}
		return nil, err
	}

	return &reception, nil
}

func (r *receptionRepo) CloseReception(reception *model.Reception) error {
	reception.Status = "close"
	return r.db.Save(reception).Error
}
