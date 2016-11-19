package main

import (
	"github.com/pressly/chi"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
