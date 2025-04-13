package main

import (
	"log"
	"net/http"

	"pvz/internal/config"
	"pvz/internal/storage/postgres"

	"pvz/internal/delivery/routes"
)

func main() {
	cfg := config.LoadConfig()
	db := postgres.InitDB(&cfg)
	postgres.Migrate(db)

	r := routes.Setup(db)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
