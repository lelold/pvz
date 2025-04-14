package postgres

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)
	if err != nil {
		log.Fatalf("failed to enable pgcrypto extension: %v", err)
	}
	queries := []string{
		`CREATE TABLE users (
    		id UUID PRIMARY KEY,
    		email TEXT UNIQUE NOT NULL,
    		password TEXT NOT NULL,
    		role TEXT NOT NULL
		);`,

		`CREATE TABLE pvzs (
    		id UUID PRIMARY KEY,
    		registration_date TIMESTAMP NOT NULL,
    		city TEXT NOT NULL
		);`,

		`CREATE TABLE receptions (
    		id UUID PRIMARY KEY,
    		date_time TIMESTAMP NOT NULL,
    		pvz_id UUID REFERENCES pvzs(id),
    		status TEXT NOT NULL
		);`,

		`CREATE TABLE products (
    		id UUID PRIMARY KEY,
    		date_time TIMESTAMP NOT NULL,
    		type TEXT NOT NULL,
    		reception_id UUID REFERENCES receptions(id)
		);`,
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	for _, q := range queries {
		_, err := tx.Exec(q)
		if err != nil {
			tx.Rollback()
			log.Fatalf("migration failed: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit migration transaction: %v", err)
	}

	log.Println("migrations are done")
}
