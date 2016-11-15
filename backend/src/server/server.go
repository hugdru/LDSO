package main

import (
	"net/http"
	"server/conn"
	"server/db"
	"github.com/pressly/chi"
	"github.com/rs/cors"
)

func main() {
	dbSession := db.Connect("mongodb:27017")
	defer db.Disconnect(dbSession)
	router := chi.NewRouter()

	db.SetCollections(dbSession, "Places4All", "main_group", "sub_group",
			"criterion", "accessibility", "property")

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)

	setRouter(router, "/mainGroups", "main_group")
	setRouter(router, "/subGroups", "sub_group")
	setRouter(router, "/criteria", "criterion")
	setRouter(router, "/accessibilities", "accessibility")

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func setRouter(router *chi.Mux, url, coll string) {
	router.Get(url, conn.Get(db.GetCollection(coll)))
	router.Get(url + "/find", conn.GetOne(db.GetCollection(coll)))
	router.Post(url, conn.Set(db.GetCollection(coll)))
	router.Put(url, conn.Update(db.GetCollection(coll)))
	router.Delete(url, conn.Delete(db.GetCollection(coll)))
}
