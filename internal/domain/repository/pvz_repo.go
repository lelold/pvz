package repository

import (
	"database/sql"
	"fmt"
	"log"
	"pvz/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type PVZRepo interface {
	Create(pvz *model.PVZ) error
	GetFilteredPVZs(start, end *time.Time, page, limit int) ([]model.PVZ, error)
	GetReceptionsWithProducts(pvzID uuid.UUID) ([]model.ReceptionWithProducts, error)
}

type pvzRepo struct {
	db *sql.DB
}

func NewPVZRepo(db *sql.DB) PVZRepo {
	return &pvzRepo{db: db}
}

func (r *pvzRepo) Create(pvz *model.PVZ) (err error) {
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

	query := `
		INSERT INTO pvzs (id, registration_date, city)
		VALUES ($1, $2, $3)
	`
	if pvz.ID == uuid.Nil {
		pvz.ID = uuid.New()
	}
	_, err = tx.Exec(query, pvz.ID, pvz.RegistrationDate, pvz.City)
	return err
}

func (r *pvzRepo) GetFilteredPVZs(start, end *time.Time, page, limit int) ([]model.PVZ, error) {
	query := `SELECT id, registration_date, city FROM pvzs WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if start != nil {
		query += " AND registration_date >= $" + itoa(argID)
		args = append(args, *start)
		argID++
	}
	if end != nil {
		query += " AND registration_date <= $" + itoa(argID)
		args = append(args, *end)
		argID++
	}

	offset := (page - 1) * limit
	query += " ORDER BY registration_date DESC LIMIT $" + itoa(argID) + " OFFSET $" + itoa(argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("Database query error: %v, query: %s", err, query)
		return nil, err
	}
	defer rows.Close()

	var pvzs []model.PVZ
	for rows.Next() {
		var pvz model.PVZ
		err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}
	return pvzs, nil
}

func (r *pvzRepo) GetReceptionsWithProducts(pvzID uuid.UUID) ([]model.ReceptionWithProducts, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	queryReceptions := `SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1`
	rows, err := tx.Query(queryReceptions, pvzID)
	if err != nil {
		log.Printf("Database query error: %v, query: %s", err, queryReceptions)
		return nil, err
	}
	defer rows.Close()

	var result []model.ReceptionWithProducts
	for rows.Next() {
		var rcp model.Reception
		err := rows.Scan(&rcp.ID, &rcp.DateTime, &rcp.PVZID, &rcp.Status)
		if err != nil {
			return nil, err
		}

		queryProducts := `SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1`
		prodRows, err := tx.Query(queryProducts, rcp.ID.String())
		log.Print(rcp.ID)
		if err != nil {
			log.Printf("Database query error: %v, query: %s", err, queryProducts)
			return nil, err
		}

		var products []model.Product
		for prodRows.Next() {
			var p model.Product
			err := prodRows.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionID)
			if err != nil {
				prodRows.Close()
				return nil, fmt.Errorf("failed to scan product: %w", err)
			}
			products = append(products, p)
		}
		prodRows.Close()

		result = append(result, model.ReceptionWithProducts{
			Reception: rcp,
			Products:  products,
		})
	}

	return result, nil
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
