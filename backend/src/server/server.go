package main

import (
	"net/http"
	"server/conn"
	"server/db"
	"github.com/pressly/chi"
)

func main() {

	dbSession := db.Connect("mongodb:27017")
	defer db.Disconnect(dbSession)

	db.SetCollection(dbSession, "Places4All", "main_group")
	db.SetCollection(dbSession, "Places4All", "sub_group")
	db.SetCollection(dbSession, "Places4All", "criterion")
	db.SetCollection(dbSession, "Places4All", "accessibility")
	db.SetCollection(dbSession, "Places4All", "property")

	router := chi.NewRouter()

	router.Get("/mainGroup", conn.Get(db.GetCollection("main_group")))
	router.Get("/mainGroup/find", conn.GetOne(db.GetCollection("main_group")))
	router.Post("/mainGroup", conn.Set(db.GetCollection("main_group")))
	router.Put("/mainGroup", conn.Update(db.GetCollection("main_group")))
	router.Delete("/mainGroup", conn.Remove(db.GetCollection("main_group")))

	router.Get("/subGroup", conn.Get(db.GetCollection("sub_group")))
	router.Get("/subGroup/find", conn.GetOne(db.GetCollection("sub_group")))
	router.Post("/subGroup", conn.Set(db.GetCollection("sub_group")))
	router.Put("/subGroup", conn.Update(db.GetCollection("sub_group")))
	router.Delete("/subGroup", conn.Remove(db.GetCollection("sub_group")))

	router.Get("/criterion", conn.Get(db.GetCollection("criterion")))
	router.Get("/criterion/find", conn.GetOne(db.GetCollection("criterion")))
	router.Post("/criterion", conn.Set(db.GetCollection("criterion")))
	router.Put("/criterion", conn.Update(db.GetCollection("criterion")))
	router.Delete("/criterion", conn.Remove(db.GetCollection("criterion")))

	router.Get("/accessibility", conn.Get(db.GetCollection("accessibility")))
	router.Get("/accessibility/find", conn.GetOne(db.GetCollection("accessibility")))
	router.Post("/accessibility", conn.Set(db.GetCollection("accessibility")))
	router.Put("/accessibility", conn.Update(db.GetCollection("accessibility")))
	router.Delete("/accessibility", conn.Remove(db.GetCollection("accessibility")))

	http.ListenAndServe(":8080", router)
}
