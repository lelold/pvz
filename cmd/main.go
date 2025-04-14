package main

import (
	"log"
	"net/http"
	"os"

	"pvz/internal/config"
	"pvz/internal/storage/postgres"

	"pvz/internal/delivery/routes"
)

func main() {
	cfg := config.LoadConfig()
	db := postgres.InitDB(&cfg)
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		postgres.Migrate(db)
		log.Println("Migration completed")
		return
	}

	r := routes.Setup(db)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
