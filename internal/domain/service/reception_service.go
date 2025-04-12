package service

import (
	"errors"
	"time"

	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"

	"github.com/google/uuid"
)

type ReceptionService interface {
	StartReception(pvzID string) (*model.Reception, error)
	CloseLastReception(pvzID string) (*model.Reception, error)
}

type receptionService struct {
	Repo repository.ReceptionRepo
}

func NewReceptionService(repo repository.ReceptionRepo) *receptionService {
	return &receptionService{Repo: repo}
}

func (s *receptionService) StartReception(pvzID string) (*model.Reception, error) {
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

func (s *receptionService) CloseLastReception(pvzID string) (*model.Reception, error) {
	reception, err := s.Repo.FindLastOpenReceptionByPVZ(pvzID)
	if err != nil {
		return nil, err
	}

	if err := s.Repo.CloseReception(reception); err != nil {
		return nil, err
	}

	return reception, nil
}
