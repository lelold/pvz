package service

import (
	"errors"
	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

type PVZService interface {
	GetFullPVZList(start, end *time.Time, page, limit int) ([]model.PVZFull, error)
	CreatePVZ(city string) (*model.PVZ, error)
}

type pvzService struct {
	repo repository.PVZRepository
}

func NewPVZService(repo repository.PVZRepository) PVZService {
	return &pvzService{repo: repo}
}

func (s *pvzService) CreatePVZ(city string) (*model.PVZ, error) {
	if !allowedCities[city] {
		return nil, errors.New("город не разрешен")
	}

	pvz := &model.PVZ{
		ID:               uuid.New(),
		City:             city,
		RegistrationDate: time.Now(),
	}
	if err := s.repo.Create(pvz); err != nil {
		return nil, err
	}

	return pvz, nil
}

func (s *pvzService) GetFullPVZList(start, end *time.Time, page, limit int) ([]model.PVZFull, error) {
	pvzs, err := s.repo.GetFilteredPVZs(start, end, page, limit)
	if err != nil {
		return nil, err
	}

	var fullList []model.PVZFull
	for _, pvz := range pvzs {
		receptions, err := s.repo.GetReceptionsWithProducts(pvz.ID.String())
		if err != nil {
			return nil, err
		}
		fullList = append(fullList, model.PVZFull{
			PVZ:        pvz,
			Receptions: receptions,
		})
	}

	return fullList, nil
}
