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

	corsm := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	dstore := datastore.Connect()
	defer dstore.Close()
	hand := &handler.Handler{Datastore: dstore}

	router.Use(dstore.SessionManager)
	router.Use(corsm.Handler)
	hand.Init(router)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
