package service

import (
	"errors"
	"time"

	"pvz/internal/domain/model"
	"pvz/internal/repository"

	"github.com/google/uuid"
)

type ReceptionService struct {
	Repo repository.ReceptionRepo
}

func NewReceptionService(repo repository.ReceptionRepo) *ReceptionService {
	return &ReceptionService{Repo: repo}
}

func (s *ReceptionService) StartReception(pvzID string) (*model.Reception, error) {
	open, err := s.Repo.HasOpenReception(pvzID)
	if err != nil {
		return nil, err
	}
	if open {
		return nil, errors.New("предыдущая приёмка ещё не закрыта")
	}

	reception := &model.Reception{
		ID:       uuid.New(),
		PVZID:    uuid.MustParse(pvzID),
		DateTime: time.Now(),
		Status:   "in_progress",
	}

	err = s.Repo.CreateReception(reception)
	if err != nil {
		return nil, err
	}

	return reception, nil
}
