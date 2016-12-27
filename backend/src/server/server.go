package main

import (
	"github.com/pressly/chi"
	"github.com/rs/cors"
	"net/http"
	"server/datastore"
	"server/handler"
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
	store := datastore.Connect()
	defer store.Close()
	hand := &handler.Handler{Datastore: store}

	hand.Init(router)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
