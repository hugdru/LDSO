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

	router.Get("/mainGroups", conn.Get(db.GetCollection("main_group")))
	router.Get("/mainGroups/find", conn.GetOne(db.GetCollection("main_group")))
	router.Post("/mainGroups", conn.Set(db.GetCollection("main_group")))
	router.Put("/mainGroups", conn.Update(db.GetCollection("main_group")))
	router.Delete("/mainGroups", conn.Delete(db.GetCollection("main_group")))

	router.Get("/subGroups", conn.Get(db.GetCollection("sub_group")))
	router.Get("/subGroups/find", conn.GetOne(db.GetCollection("sub_group")))
	router.Post("/subGroups", conn.Set(db.GetCollection("sub_group")))
	router.Put("/subGroups", conn.Update(db.GetCollection("sub_group")))
	router.Delete("/subGroups", conn.Delete(db.GetCollection("sub_group")))

	router.Get("/criteria", conn.Get(db.GetCollection("criterion")))
	router.Get("/criteria/find", conn.GetOne(db.GetCollection("criterion")))
	router.Post("/criteria", conn.Set(db.GetCollection("criterion")))
	router.Put("/criteria", conn.Update(db.GetCollection("criterion")))
	router.Delete("/criteria", conn.Delete(db.GetCollection("criterion")))

	router.Get("/accessibilities", conn.Get(db.GetCollection("accessibility")))
	router.Get("/accessibilities/find", conn.GetOne(db.GetCollection("accessibility")))
	router.Post("/accessibilities", conn.Set(db.GetCollection("accessibility")))
	router.Put("/accessibilities", conn.Update(db.GetCollection("accessibility")))
	router.Delete("/accessibilities", conn.Delete(db.GetCollection("accessibility")))

	http.ListenAndServe(":8080", router)
}
