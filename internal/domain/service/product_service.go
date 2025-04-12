package service

import (
	"errors"
	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

var allowedProducts = map[string]bool{
	"одежда":      true,
	"обувь":       true,
	"электроника": true,
}

type ProductService interface {
	CreateProduct(typeStr string, pvzId string) (*model.Product, error)
	DeleteLastProduct(pvzID, userRole string) error
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
	if !allowedProducts[typeStr] {
		return nil, errors.New("неверный товар")
	}
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

func (s *productService) DeleteLastProduct(pvzID, userRole string) error {
	if userRole != "employee" {
		return errors.New("доступ запрещён")
	}
	reception, err := s.receptionRepo.FindLastOpenReceptionByPVZ(pvzID)
	if err != nil {
		return errors.New("не удалось найти последнюю открытую приемку")
	}
	product, err := s.productRepo.GetLastAddedProduct(reception.ID.String())
	if err != nil {
		return errors.New("не удалось найти последний добавленный товар")
	}
	err = s.productRepo.DeleteProductByID(product.ID.String())
	if err != nil {
		return errors.New("не удалось удалить товар")
	}
	return nil
}
