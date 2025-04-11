package service

import (
	"errors"
	"pvz/internal/domain/model"
	"pvz/internal/repository"
	"time"

	"github.com/google/uuid"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

type PVZService struct {
	repo repository.PVZRepo
}

func NewPVZService(repo repository.PVZRepo) *PVZService {
	return &PVZService{repo: repo}
}

func (s *PVZService) CreatePVZ(city string) (*model.PVZ, error) {
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

func (s *PVZService) GetAllPVZ() ([]model.PVZ, error) {
	return s.repo.GetAll()
}
