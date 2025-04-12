package main

import (
	"log"
	"net/http"

	"pvz/internal/config"
	"pvz/internal/delivery/handler"
	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/service"
	"pvz/internal/postgres"
	"pvz/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db := postgres.InitDB(&cfg)
	postgres.Migrate(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	receptionRepo := repository.NewReceptionRepo(db)
	receptionService := service.NewReceptionService(receptionRepo)
	receptionHandler := handler.NewReceptionHandler(receptionService)

	pvzRepo := repository.NewPVZRepo(db)
	pvzService := service.NewPVZService(pvzRepo)
	pvzHandler := handler.NewPVZHandler(pvzService)

	authHandler := &handler.AuthHandler{UserService: userService}

	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register)
	r.HandleFunc("/login", authHandler.Login)
	r.HandleFunc("/dummyLogin", handler.DummyLoginHandler)
	r.Handle("/pvz", middleware.AuthMiddleware(http.HandlerFunc(pvzHandler.HandlePVZ)))
	r.Handle("/receptions", middleware.AuthMiddleware(http.HandlerFunc(receptionHandler.StartReception))).Methods("POST")
	r.Handle("/pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(http.HandlerFunc(receptionHandler.CloseLastReception))).Methods("POST")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
