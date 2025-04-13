package postgres

import (
	"log"
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.PVZ{},
		&model.Reception{},
		&model.Product{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations are done")
}
