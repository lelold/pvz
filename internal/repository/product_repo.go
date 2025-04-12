package repository

import (
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

type ProductRepo interface {
	Create(product *model.Product) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product *model.Product) error {
	return r.db.Create(product).Error
}
