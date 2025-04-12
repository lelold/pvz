package main

import (
	"log"
	"net/http"

	"pvz/internal/config"
	"pvz/internal/delivery/handler"
	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/repository"
	"pvz/internal/domain/service"
	"pvz/internal/postgres"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db := postgres.InitDB(&cfg)
	postgres.Migrate(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := &handler.AuthHandler{UserService: userService}

	receptionRepo := repository.NewReceptionRepo(db)
	receptionService := service.NewReceptionService(receptionRepo)
	receptionHandler := handler.NewReceptionHandler(receptionService)

	pvzRepo := repository.NewPVZRepo(db)
	pvzService := service.NewPVZService(pvzRepo)
	pvzHandler := handler.NewPVZHandler(pvzService)

	productRepo := repository.NewProductRepo(db)
	productService := service.NewProductService(productRepo, receptionRepo)
	productHandler := handler.NewProductHandler(productService)

	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/dummyLogin", handler.DummyLoginHandler).Methods("POST")
	r.Handle("/pvz", middleware.AuthMiddleware(http.HandlerFunc(pvzHandler.HandlePVZ)))
	r.Handle("/receptions", middleware.AuthMiddleware(http.HandlerFunc(receptionHandler.StartReception))).Methods("POST")
	r.Handle("/pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(http.HandlerFunc(receptionHandler.CloseLastReception))).Methods("POST")
	r.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(productHandler.CreateProduct))).Methods("POST")
	r.Handle("/pvz/{pvzId}/delete_last_product", middleware.AuthMiddleware(http.HandlerFunc(productHandler.DeleteLastProduct))).Methods("POST")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
