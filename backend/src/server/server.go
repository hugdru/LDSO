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
	router := chi.NewRouter()

	db.SetCollections(dbSession, "Places4All", "main_group", "sub_group",
			"criterion", "accessibility", "property")

	setRouter(router, "/mainGroups", "main_group")
	setRouter(router, "/subGroups", "sub_group")
	setRouter(router, "/criteria", "criterion")
	setRouter(router, "/accessibilities", "accessibility")

	http.ListenAndServe(":8080", router)
}

func setRouter(router *chi.Mux, url, coll string) {
	router.Get(url, conn.Get(db.GetCollection(coll)))
	router.Get(url + "/find", conn.GetOne(db.GetCollection(coll)))
	router.Post(url, conn.Set(db.GetCollection(coll)))
	router.Put(url, conn.Update(db.GetCollection(coll)))
	router.Delete(url, conn.Delete(db.GetCollection(coll)))

}
