package main

import (
	"log"
	"net/http"

	"pvz/internal/config"
	"pvz/internal/delivery/pvz_http"
	"pvz/internal/domain/service"
	"pvz/internal/postgres"
	repository "pvz/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	db := postgres.InitDB(&cfg)
	postgres.Migrate(db)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := &pvz_http.AuthHandler{UserService: userService}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/dummyLogin", pvz_http.DummyLoginHandler)
	mux.Handle("/pvz", pvz_http.AuthMiddleware(pvz_http.PVZHandler()))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
