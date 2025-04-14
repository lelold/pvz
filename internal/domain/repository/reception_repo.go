package repository

import (
	"database/sql"
	"errors"
	"pvz/internal/domain/model"

	"github.com/google/uuid"
)

type ReceptionRepo interface {
	CreateReception(reception *model.Reception) error
	FindLastOpenReceptionByPVZ(pvzID uuid.UUID) (*model.Reception, error)
	CloseReception(reception *model.Reception) error
	HasOpenReception(pvzID uuid.UUID) (bool, error)
}

type receptionRepo struct {
	db *sql.DB
}

func NewReceptionRepo(db *sql.DB) ReceptionRepo {
	return &receptionRepo{db: db}
}

func (r *receptionRepo) CreateReception(reception *model.Reception) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if reception.ID == uuid.Nil {
		reception.ID = uuid.New()
	}
	query := `
		INSERT INTO receptions (id, date_time, pvz_id, status)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(query, reception.ID, reception.DateTime, reception.PVZID, reception.Status)
	return err
}

func (r *receptionRepo) HasOpenReception(pvzID uuid.UUID) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress'
	`
	err := r.db.QueryRow(query, pvzID).Scan(&count)
	return count > 0, err
}

func (r *receptionRepo) FindLastOpenReceptionByPVZ(pvzID uuid.UUID) (*model.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status 
		FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`
	row := r.db.QueryRow(query, pvzID)
	var rcp model.Reception
	err := row.Scan(&rcp.ID, &rcp.DateTime, &rcp.PVZID, &rcp.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("открытая приемка не найдена")
		}
		return nil, err
	}
	return &rcp, nil
}

func (r *receptionRepo) CloseReception(reception *model.Reception) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `UPDATE receptions SET status = 'close' WHERE id = $1`
	_, err = tx.Exec(query, reception.ID)
	if err == nil {
		reception.Status = "close"
	}
	return err
}
