package main

import (
	"log"
	"net/http"

	"pvz/internal/delivery/pvz_http"
	"pvz/internal/domain/service"
	repository "pvz/internal/repository/pg"
)

func main() {
	userRepo := repository.NewUserRepo()
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
