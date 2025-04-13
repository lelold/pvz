package postgres

import (
	"fmt"
	"log"
	"pvz/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(conf *config.Config) *gorm.DB {
	var DB *gorm.DB
	conn_str := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
		conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName,
	)

	conn, err := gorm.Open(postgres.Open(conn_str), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB: ", err)
	}
	DB = conn

	log.Println("DB succesfully connected")
	return DB
}
