package repository

import (
	"database/sql"
	"pvz/internal/domain/model"

	"github.com/google/uuid"
)

type ProductRepo interface {
	Create(product *model.Product) error
	GetLastAddedProduct(receptionID uuid.UUID) (*model.Product, error)
	DeleteProductByID(productID uuid.UUID) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product *model.Product) error {
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
		INSERT INTO products (id, date_time, type, reception_id)
		VALUES ($1, $2, $3, $4);
	`
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	_, err = tx.Exec(query, product.ID, product.DateTime, product.Type, product.ReceptionID)
	return err
}

func (r *productRepo) GetLastAddedProduct(receptionID uuid.UUID) (*model.Product, error) {
	query := `
		SELECT id, date_time, type, reception_id
		FROM products
		WHERE reception_id = $1
		ORDER BY date_time DESC
		LIMIT 1;
	`

	var product model.Product
	err := r.db.QueryRow(query, receptionID).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) DeleteProductByID(productID uuid.UUID) error {
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

	query := `DELETE FROM products WHERE id = $1;`
	_, err = tx.Exec(query, productID)
	return err
}
