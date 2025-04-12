package service

import (
	"errors"
	"pvz/internal/domain/model"
	"pvz/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(typeStr string, pvzId string) (*model.Product, error)
}

type productService struct {
	productRepo   repository.ProductRepo
	receptionRepo repository.ReceptionRepo
}

func NewProductService(pRepo repository.ProductRepo, rRepo repository.ReceptionRepo) ProductService {
	return &productService{
		productRepo:   pRepo,
		receptionRepo: rRepo,
	}
}

func (s *productService) CreateProduct(typeStr string, pvzId string) (*model.Product, error) {
	reception, err := s.receptionRepo.FindLastOpenReceptionByPVZ(pvzId)
	if err != nil {
		return nil, errors.New("нет активной приёмки")
	}

	product := &model.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        typeStr,
		ReceptionID: reception.ID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}
