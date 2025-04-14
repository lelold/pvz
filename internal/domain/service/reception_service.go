package service

import (
	"errors"
	"log"
	"time"

	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"

	"github.com/google/uuid"
)

type ReceptionService interface {
	StartReception(pvzID uuid.UUID) (*model.Reception, error)
	CloseLastReception(pvzID uuid.UUID) (*model.Reception, error)
}

type receptionService struct {
	Repo repository.ReceptionRepo
}

func NewReceptionService(repo repository.ReceptionRepo) *receptionService {
	return &receptionService{Repo: repo}
}

func (s *receptionService) StartReception(pvzID uuid.UUID) (*model.Reception, error) {
	open, err := s.Repo.HasOpenReception(pvzID)
	if err != nil {
		return nil, err
	}
	if open {
		return nil, errors.New("предыдущая приёмка ещё не закрыта")
	}

	reception := &model.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		DateTime: time.Now(),
		Status:   "in_progress",
	}
	log.Printf(reception.ID.String(), reception.PVZID.String())
	err = s.Repo.CreateReception(reception)
	if err != nil {
		return nil, err
	}

	return reception, nil
}

func (s *receptionService) CloseLastReception(pvzID uuid.UUID) (*model.Reception, error) {
	reception, err := s.Repo.FindLastOpenReceptionByPVZ(pvzID)
	if err != nil {
		return nil, err
	}

	if err := s.Repo.CloseReception(reception); err != nil {
		return nil, err
	}

	return reception, nil
}
