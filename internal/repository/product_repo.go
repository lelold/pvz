package repository

import (
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

type ProductRepo interface {
	Create(product *model.Product) error
	GetLastAddedProduct(receptionID string) (*model.Product, error)
	DeleteProductByID(productID string) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepo) GetLastAddedProduct(receptionID string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("reception_id = ?", receptionID).
		Order("date_time DESC").
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) DeleteProductByID(productID string) error {
	return r.db.Delete(&model.Product{}, "id = ?", productID).Error
}
