package repository

import (
	"pvz/internal/domain/model"
	"time"

	"gorm.io/gorm"
)

type PVZRepository interface {
	Create(pvz *model.PVZ) error
	GetFilteredPVZs(start, end *time.Time, page, limit int) ([]model.PVZ, error)
	GetReceptionsWithProducts(pvzID string) ([]model.ReceptionWithProducts, error)
}

type pvzRepo struct {
	db *gorm.DB
}

func NewPVZRepo(db *gorm.DB) PVZRepository {
	return &pvzRepo{db: db}
}

func (r *pvzRepo) Create(pvz *model.PVZ) error {
	return r.db.Create(pvz).Error
}

func (r *pvzRepo) GetFilteredPVZs(start, end *time.Time, page, limit int) ([]model.PVZ, error) {
	var pvzs []model.PVZ
	query := r.db.Model(&model.PVZ{})

	if start != nil {
		query = query.Where("registration_date >= ?", *start)
	}
	if end != nil {
		query = query.Where("registration_date <= ?", *end)
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&pvzs).Error; err != nil {
		return nil, err
	}
	return pvzs, nil
}

func (r *pvzRepo) GetReceptionsWithProducts(pvzID string) ([]model.ReceptionWithProducts, error) {
	var receptions []model.Reception
	if err := r.db.Where("pvz_id = ?", pvzID).Find(&receptions).Error; err != nil {
		return nil, err
	}

	var result []model.ReceptionWithProducts
	for _, rcp := range receptions {
		var products []model.Product
		if err := r.db.Where("reception_id = ?", rcp.ID).Find(&products).Error; err != nil {
			return nil, err
		}
		result = append(result, model.ReceptionWithProducts{
			Reception: rcp,
			Products:  products,
		})
	}
	return result, nil
}
