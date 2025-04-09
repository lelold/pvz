package main

import (
	"log"
	"net/http"

	"pvz/internal/delivery/pvz_http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/dummyLogin", pvz_http.DummyLoginHandler)
	mux.Handle("/pvz", pvz_http.AuthMiddleware(pvz_http.PVZHandler()))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
